package redisfunctions

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func ConnectRedisDependency(t *testing.T) *redis.Client {
	t.Helper()

	rdb := ConnectRedis()

	t.Cleanup(func() {
		rdb.Close()
	})

	return rdb
}

func TestPubsubReverser(t *testing.T) {
	ctx := context.Background()
	rdb := ConnectRedisDependency(t)

	testCases := []struct {
		sourceChannel      string
		destinationChannel string
		durationMs         int
		messages           []string
		expectedResponses  []string
	}{
		{
			"src1",
			"dest1",
			250,
			[]string{"ping", "pong", "king", "kong"},
			[]string{"gnip", "gnop"},
		},
		{
			"src2",
			"dest2",
			450,
			[]string{"public", "static", "final", "void", "main"},
			[]string{"cilbup", "citats", "lanif", "diov"},
		},
	}

	for _, testCase := range testCases {
		responses := PubsubReverserSingleTest(
			t,
			ctx,
			rdb,
			testCase.sourceChannel,
			testCase.destinationChannel,
			testCase.durationMs,
			&(testCase.messages),
		)

		assert.Equal(t, testCase.expectedResponses, responses)
	}
}

func PubsubReverserSingleTest(
	t *testing.T,
	ctx context.Context,
	rdb *redis.Client,
	sourceChannel string,
	destinationChannel string,
	durationMs int,
	messages *[]string,
) []string {
	responses := []string{}

	pubsub := rdb.Subscribe(ctx, destinationChannel)
	defer pubsub.Close()
	pubsubChannel := pubsub.Channel()

	ticker := time.NewTicker(100 * time.Millisecond)
	timeout := time.NewTimer(time.Duration(1.5*float64(durationMs)) * time.Millisecond)
	timeIsOver := false

	go PubsubReverser(ctx, rdb, sourceChannel, destinationChannel, durationMs)
	go func() {
		for _, message := range *messages {
			<-ticker.C
			rdb.Publish(ctx, sourceChannel, message)
		}
	}()

	for !timeIsOver {
		select {
		case message := <-pubsubChannel:
			responses = append(responses, message.Payload)
		case <-timeout.C:
			timeIsOver = true
		}
	}

	return responses
}