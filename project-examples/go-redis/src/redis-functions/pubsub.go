package redisfunctions

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func PubsubReverser(
	ctx context.Context,
	rdb *redis.Client,
	sourceChannel string,
	destinationChannel string,
	ms int,
) {
	pubsub := rdb.Subscribe(ctx, sourceChannel)
	defer pubsub.Close()
	pubsubChannel := pubsub.Channel()

	timeout := time.NewTimer(time.Duration(ms) * time.Millisecond)
	timeIsOver := false

	for !timeIsOver {
		select {
		case message := <-pubsubChannel:
			response := ReverseString(message.Payload)
			err := rdb.Publish(ctx, destinationChannel, response).Err()
			if err != nil {
				panic(err)
			}
		case <-timeout.C:
			timeIsOver = true
		}
	}
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}