package internal

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/reconmap/shared-lib/pkg/api"
	reconmapio "github.com/reconmap/shared-lib/pkg/io"
	"github.com/reconmap/shared-lib/pkg/logging"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

// App contains properties needed for agent
// to connect to redis and http router.
type App struct {
	redisConn *redis.Client
	muxRouter *mux.Router
	Logger    *zap.SugaredLogger

	debugEnabled  bool
	skipTlsVerify bool
}

var logger = logging.GetLoggerInstance()

// NewApp returns a App struct that has intialized a redis client and http router.
func NewApp() App {
	debugValue, _ := os.LookupEnv("RMAP_KEYCLOAK_DEBUG")
	skipTlsVerify, _ := os.LookupEnv("RMAP_KEYCLOAK_SKIP_TLS_VERIFY")

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/term", handleWebsocket)
	muxRouter.HandleFunc("/notifications", handleNotifications)

	return App{
		muxRouter:     muxRouter,
		Logger:        logging.GetLoggerInstance(),
		debugEnabled:  debugValue == "true",
		skipTlsVerify: skipTlsVerify == "true",
	}
}

// Run starts the agent.
func (app *App) Run() *error {
	app.Logger.Info("Reconmap agent starting...")

	accessToken, err := GetAccessToken(app)
	if err != nil {
		app.Logger.Error("unable to login to keycloak", zap.Error(err))
		panic(err)
	}

	schedules, err := api.GetCommandsSchedules(accessToken)
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
				stdout, errStdout = reconmapio.CopyAndCapture(os.Stdout, stdoutIn)
				wg.Done()
			}()

			stderr, errStderr = reconmapio.CopyAndCapture(os.Stderr, stderrIn)

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
		return &errorFormatted
	}

	go broadcastNotifications(app)

	if err := http.ListenAndServe(*listen, app.muxRouter); err != nil {
		app.Logger.Fatal("Something went wrong with the webserver", zap.Error(err))
	}

	return nil
}
