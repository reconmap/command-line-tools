package internal

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"reconmap/agent/internal/build"
	"reconmap/agent/internal/configuration"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/reconmap/shared-lib/pkg/api"
	sharedconfig "github.com/reconmap/shared-lib/pkg/configuration"
	sharedio "github.com/reconmap/shared-lib/pkg/io"
	"github.com/reconmap/shared-lib/pkg/logging"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
)

// App contains properties needed for agent
// to connect to redis and http router.
type App struct {
	redisConn *redis.Client
	muxRouter *mux.Router
	Logger    *zap.SugaredLogger
}

var logger = logging.GetLoggerInstance()

// NewApp returns a App struct that has intialized a redis client and http router.
func NewApp() App {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/term", handleWebsocket)
	muxRouter.HandleFunc("/notifications", handleNotifications)
	muxRouter.HandleFunc("/systeminfo", handleSystemInfo)

	return App{
		muxRouter: muxRouter,
		Logger:    logging.GetLoggerInstance(),
	}
}

// Run starts the agent.
func (app *App) Run() error {
	app.Logger.Info("Reconmap agent starting...")

	config, err := sharedconfig.ReadConfig[configuration.Config]("config-reconmapd.json")
	if err != nil {
		app.Logger.Error("unable to read reconmapd config", zap.Error(err))
		return err
	}

	accessToken, err := GetAccessToken(app)
	if err != nil {
		return fmt.Errorf("unable to login to keycloak (%w)", err)
	}

	restApiUrl := config.ReconmapApiConfig.BaseUri
	schedules, err := api.GetCommandsSchedules(restApiUrl, accessToken)
	if err != nil {
		app.Logger.Error("unable to get command schedules", zap.Error(err))
	}

	app.Logger.Info("creating cron jobs")
	c := cron.New()

	for _, commandSchedule := range *schedules {
		c.AddFunc(commandSchedule.CronExpression, func() {
			parts := strings.Split(commandSchedule.ArgumentValues, " ")
			cmd := exec.Command(parts[0], parts[1:]...) // #nosec G204
			cmd.Env = append(os.Environ(), "PS1=# ")
			cmd.Env = append(cmd.Env, "TERM=xterm")
			cmd.Env = append(cmd.Env, "RMAP_SESSION_TOKEN="+accessToken)
			var stdout, stderr []byte
			var errStdout, errStderr error
			stdoutIn, _ := cmd.StdoutPipe()
			stderrIn, _ := cmd.StderrPipe()
			err := cmd.Start()
			if err != nil {
				app.Logger.Fatalf("cmd.Start() failed with '%s'\n", err)
			}
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				stdout, errStdout = sharedio.CopyAndCapture(os.Stdout, stdoutIn)
				wg.Done()
			}()

			stderr, errStderr = sharedio.CopyAndCapture(os.Stderr, stderrIn)

			wg.Wait()

			err = cmd.Wait()
			if err != nil {
				if errStderr != nil {
					print(errStderr)
				}
				app.Logger.Fatalf("cmd.Run() failed with %s\n", err)
			}
			if errStdout != nil || errStderr != nil {
				app.Logger.Fatal("failed to capture stdout or stderr\n")
			}
			outStr, errStr := string(stdout), string(stderr)
			app.Logger.Debug(outStr)
			app.Logger.Debug(errStr)
		})
	}
	c.Start()

	listen := flag.String("listen", ":5520", "Host:port to listen on")
	flag.Parse()

	redisErr := app.connectRedis()
	if redisErr != nil {
		errorFormatted := fmt.Errorf("unable to connect to redis (%w)", *redisErr)
		return errorFormatted
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go func() {

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		cpuInfo, err := cpu.Counts(true)
		cpuStr := "unknown"

		if err == nil {
			cpuStr = fmt.Sprintf("%d cores", cpuInfo)
		}

		vmStat, err := mem.VirtualMemory()
		memStr := "unknown"
		if err == nil {
			memStr = fmt.Sprintf("%.2fGB", float64(vmStat.Total)/(1024*1024*1024))
		}

		systemInfo := api.SystemInfo{
			Version:  build.BuildVersion,
			Hostname: hostname,
			Arch:     runtime.GOARCH,
			Os:       runtime.GOOS,
			CPU:      cpuStr,
			Memory:   memStr,
		}

		api.AgentBoot(restApiUrl, accessToken, &systemInfo)
	}()
	go func() {
		for range ticker.C {
			app.Logger.Info("pinging reconmap API")
			api.AgentPing(restApiUrl, accessToken)
		}
	}()

	go broadcastNotifications(app)

	if err := http.ListenAndServe(*listen, app.muxRouter); err != nil {
		app.Logger.Fatal("Something went wrong with the webserver", zap.Error(err))
	}

	return nil
}
