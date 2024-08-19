package internal

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	redisTimeout = 2 * time.Second
)

func broadcastNotifications(app *App) {
	for {
		log.Debug("searching for notifications...")
		ctx := context.Background()
		result, err := app.redisConn.BRPop(ctx, redisTimeout, "notifications:queue").Result()
		if err != nil {
			log.Debug("no items retrieved from notifications queue: ", err)
		} else if result != nil {
			broadcast(result[1])
		}
	}
}
