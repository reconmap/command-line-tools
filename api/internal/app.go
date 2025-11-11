package internal

import (
	"fmt"
	"net/http"
	"reconmap/api/internal/configuration"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	sharedconfig "github.com/reconmap/shared-lib/pkg/configuration"
	"github.com/reconmap/shared-lib/pkg/logging"
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
	muxRouter.HandleFunc("/notifications", handleNotifications)

	return App{
		muxRouter: muxRouter,
		Logger:    logging.GetLoggerInstance(),
	}
}

// Run starts the agent.
func (app *App) Run() error {
	app.Logger.Info("Reconmap API starting...")

	_, err := sharedconfig.ReadConfig[configuration.Config]("config-reconmapd.json")
	if err != nil {
		app.Logger.Error("unable to read reconmapd config", zap.Error(err))
		return err
	}

	app.Logger.Info("creating cron jobs")

	redisErr := app.connectRedis()
	if redisErr != nil {
		errorFormatted := fmt.Errorf("unable to connect to redis (%w)", *redisErr)
		return errorFormatted
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go broadcastNotifications(app)

	if err := http.ListenAndServe(*listen, app.muxRouter); err != nil {
		app.Logger.Fatal("Something went wrong with the webserver", zap.Error(err))
	}

	return nil
}
