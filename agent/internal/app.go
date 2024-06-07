package internal

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/Nerzal/gocloak/v11"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/reconmap/shared-lib/pkg/io"
	"github.com/reconmap/shared-lib/pkg/models"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// App contains properties needed for agent
// to connect to redis and http router.
type App struct {
	redisConn *redis.Client
	muxRouter *mux.Router
}

// NewApp returns a App struct that has intialized a redis client and http router.
func NewApp() App {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/term", handleWebsocket)
	muxRouter.HandleFunc("/notifications", handleNotifications)

	return App{
		muxRouter: muxRouter,
	}
}

// Run starts the agent.
func (app *App) Run() *error {
	log.Info("Reconmap agent")

	keycloakHostname, _ := os.LookupEnv("RMAP_KEYCLOAK_HOSTNAME")
	clientID, _ := os.LookupEnv("RMAP_AGENT_CLIENT_ID")
	clientSecret, _ := os.LookupEnv("RMAP_AGENT_CLIENT_SECRET")
	restApiUrl, _ := os.LookupEnv("RMAP_REST_API_URL")

	realm := "reconmap"
	client := gocloak.NewClient(keycloakHostname, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	restyClient := client.RestyClient()
	restyClient.SetDebug(true)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	ctx := context.Background()
	token, err := client.LoginClient(ctx, clientID, clientSecret, realm)
	if err != nil {
		panic("Login failed:" + err.Error())
	}
	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, clientID, clientSecret, realm)
	if err != nil {
		panic("Inspection failed:" + err.Error())
	}

	if !*rptResult.Active {
		panic("Token is not active")
	}

	client2 := &http.Client{}
	req, err := http.NewRequest("GET", restApiUrl+"/commands/schedules", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	response, err := client2.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var schedules *models.CommandSchedules = &models.CommandSchedules{}

	if err = json.Unmarshal(body, schedules); err != nil {
		panic(err)
	}

	log.Info("Create new cron")
	c := cron.New()

	for _, commandSchedule := range *schedules {
		c.AddFunc(commandSchedule.CronExpression, func() {
			parts := strings.Split(commandSchedule.ArgumentValues, " ")
			cmd := exec.Command(parts[0], parts[1:]...) // #nosec G204
			cmd.Env = append(os.Environ(), "PS1=# ")
			cmd.Env = append(cmd.Env, "TERM=xterm")
			cmd.Env = append(cmd.Env, "RMAP_SESSION_TOKEN="+token.AccessToken)
			var stdout, stderr []byte
			var errStdout, errStderr error
			stdoutIn, _ := cmd.StdoutPipe()
			stderrIn, _ := cmd.StderrPipe()
			err := cmd.Start()
			if err != nil {
				log.Fatalf("cmd.Start() failed with '%s'\n", err)
			}
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				stdout, errStdout = io.CopyAndCapture(os.Stdout, stdoutIn)
				wg.Done()
			}()

			stderr, errStderr = io.CopyAndCapture(os.Stderr, stderrIn)

			wg.Wait()

			err = cmd.Wait()
			if err != nil {
				if errStderr != nil {
					print(errStderr)
				}
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
			if errStdout != nil || errStderr != nil {
				log.Fatal("failed to capture stdout or stderr\n")
			}
			outStr, errStr := string(stdout), string(stderr)
			log.Debug(outStr)
			log.Debug(errStr)
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
		log.WithError(err).Fatal("Something went wrong with the webserver")
	}

	return nil
}
