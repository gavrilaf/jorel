package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/testdata"
)

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	publisher, err := msgqueue.NewPublisher(ctx, testdata.ProjectID, testdata.IngressTopic)
	if err != nil {
		logger.Panicf("failed to create publisher, %v", err)
	}

	defer publisher.Close()

	publisherID := uuid.New().String()

	delays := []int{
		0,
		5,
		10,
		60,
	}

	for indx, d := range delays {
		id := fmt.Sprintf("%s-%d", publisherID, indx)
		now  := time.Now().UTC()

		m := testdata.Message{
			ID:               id,
			Created:          now,
			ScheduleDuration: d,
		}

		data, err := json.Marshal(&m)
		if err != nil {
			logger.Panicf("failed to marshal message, %v", err)
		}

		attributes := msgqueue.MsgAttributes{
			DelayInSeconds: d,
			Original:       map[string]string{"one": "two"},
		}

		res, err := publisher.Publish(ctx, data, attributes.GetAttributes())
		if err != nil {
			logger.Panicf("failed to publish message, %v", err)
		}

		res.GetMessageID(ctx)

		scheduledTime := now.Add(time.Duration(d) * time.Second)
		logger.Infof("Published message with ID: %s, duration %d, should executed in: %v", id, d, scheduledTime)
	}
}
