package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func (app *App) connectRedis() *error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	ctx := context.Background()

	redisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if _, err := redisConn.Ping(ctx).Result(); err != nil {
		return &err
	}

	app.redisConn = redisConn

	return nil
}
