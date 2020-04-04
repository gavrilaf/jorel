package main

import (
	"context"
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

	durations := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		time.Minute,
	}

	for indx, d := range durations {
		id := fmt.Sprintf("%s-%d", publisherID, indx)
		m := testdata.Message{
			ID:               id,
			Created:          time.Now(),
			ScheduleDuration: d,
		}

		res, err := publisher.Publish(ctx, m, map[string]string{})
		if err != nil {
			logger.Panicf("failed to publish message, %v", err)
		}

		messageID, err := res.GetMessageID(ctx)
		if err != nil {
			logger.Panicf("failed to read message ID, %v", err)
		}

		logger.Infof("Published message with ID: %s", messageID)
	}

}
