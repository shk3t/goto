package main

import (
	"context"
	"fmt"
	rfuncs "go-redis/src/redis-functions"
)

func main() {
	ctx := context.Background()
	rdb := rfuncs.ConnectRedis()
	defer rdb.Close()

	rfuncs.PubsubReverser(ctx, rdb, "source", "destination", 1000)
	fmt.Println("END")
}