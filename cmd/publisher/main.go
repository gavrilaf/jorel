package main

import (
	"context"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/testdata"
)

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	publisher, err := msgqueue.NewPublisher(ctx, testdata.ProjectID, testdata.TopicName)
	if err != nil {
		logger.Panicf("failed to create publisher, %v", err)
	}

	defer publisher.Close()

	m := testdata.Message{
		ID:   1,
		Text: "test",
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