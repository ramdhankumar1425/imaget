package infra

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	redisUrl := os.Getenv("REDIS_URL")

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	RDB = redis.NewClient(opt)

	if err := RDB.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
