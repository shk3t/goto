package redisfunctions

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	rdbhost := os.Getenv("RDBHOST")
	if rdbhost == "" {
		rdbhost = "localhost"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbhost + ":6379",
		Password: "",
		DB:       0,
	})

	return rdb
}