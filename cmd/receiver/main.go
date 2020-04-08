package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/testdata"
	"github.com/gavrilaf/dyson/pkg/utils"
)

const (
	projectID      = "dyson-272914"
	defaultSubscription = "default-topic-subs"
)

var (
	receivedCount = 0
	outboundCount = 0
	meanDeviation = time.Duration(0)
	maxDeviation  = time.Duration(0)
)

type handler struct{}

func (handler) Receive(ctx context.Context, data []byte, attributes map[string]string) error {
	var msg testdata.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	checked, diff := msg.Check()

	receivedCount += 1
	if !checked {
		dlog.FromContext(ctx).Warnf("** %v: Message %v, checked=%v, difference=%v", now, msg, checked, diff)
		outboundCount += 1
	} else {
		dlog.FromContext(ctx).Infof("** %v: Message %v, checked=%v, difference=%v", now, msg, checked, diff)
	}

	if diff > maxDeviation {
		maxDeviation = diff
	}

	meanDeviation += diff

	return nil
}

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	var subscriptionID = defaultSubscription
	if len(os.Args) > 1 {
		subscriptionID = os.Args[1]
	}

	receiver, err := msgqueue.NewReceiver(ctx, projectID, subscriptionID)
	if err != nil {
		logger.Panicf("failed to create receiver, %v", err)
	}

	logger.Info("Starting receiver")
	err = receiver.Run(ctx, handler{})

	utils.WaitForShutdown(ctx, func() {
		err := receiver.Close()
		logger.Infof("receiver closed, error=%v", err)

		meanDeviation = meanDeviation / time.Duration(receivedCount)

		logger.Infof("received messages %d, outbound %d, max deviation %v, mean deviation %v", receivedCount, outboundCount, maxDeviation, meanDeviation)
	})
}
