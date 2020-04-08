package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/testdata"
	"github.com/gavrilaf/dyson/pkg/utils"
)

const (
	projectID  = "dyson-272914"
	topicID = "ingress-topic"
)

const (
	simpleMode = iota
	routingMode
)

func main() {
	mode := simpleMode
	if len(os.Args) > 1 {
		mode = routingMode
	}

	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	publisher, err := msgqueue.NewPublisher(ctx, projectID, topicID)
	if err != nil {
		logger.Panicf("failed to create publisher, %v", err)
	}

	defer publisher.Close()

	publisherID := uuid.New().String()

	msgTypes := []string{"", "cancel"}

	delays := []int{
		85,
		800,
		45,
		60,
		45,
		200,
		150,
		45,
		30,
		600,
	}

	sentCount := 0
	startTime := time.Now()

	for repeat := 0; repeat < 300; repeat++ {
		for indx, d := range delays {
			id := fmt.Sprintf("%s-%d", publisherID, indx)
			now := time.Now().UTC()

			d += rand.Intn(10)

			m := testdata.Message{
				ID:               id,
				Created:          now,
				ScheduleDuration: d,
			}

			data, err := json.Marshal(&m)
			if err != nil {
				logger.Panicf("failed to marshal message, %v", err)
			}

			var attributes msgqueue.MsgAttributes

			if mode == simpleMode {
				attributes = msgqueue.MsgAttributes{
					DelayInSeconds: d,
					Original:       map[string]string{"one": "two"},
				}
			} else {
				attributes = msgqueue.MsgAttributes{
					DelayInSeconds: d,
					MessageType: msgTypes[rand.Intn(2)],
					Original:       map[string]string{"one": "two"},
				}
			}
			_, err = publisher.Publish(ctx, data, attributes.GetAttributes())
			if err != nil {
				logger.Panicf("failed to publish message, %v", err)
			}

			scheduledTime := now.Add(time.Duration(d) * time.Second)
			logger.Infof("Published message with ID: %s, duration %d, should executed in: %v", id, d, scheduledTime)

			sentCount += 1
		}

		time.Sleep(1 * time.Second)
	}

	endTime := time.Now()

	logger.Infof("Sent %d messages, %v", sentCount, endTime.Sub(startTime))

	utils.WaitForShutdown(ctx, func() {})
}
