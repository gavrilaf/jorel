package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/testdata"
)

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	receiver, err := msgqueue.NewReceiver(ctx, testdata.ProjectID, testdata.CheckSubscription)
	if err != nil {
		logger.Panicf("failed to create publisher, %v", err)
	}

	logger.Info("Starting receiver")
	err = receiver.Run(ctx, func(ctx context.Context, data []byte, attributes map[string]string) error {
		var msg testdata.Message
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return err
		}

		dlog.FromContext(ctx).Infof("Received message: %v", msg)

		checked, diff := msg.Check()
		dlog.FromContext(ctx).Infof("Message ID: %s, checked=%v, difference=%v", msg.ID, checked, diff)

		return nil
	})

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info("Ctrl+C pressed in Terminal")

		err := receiver.Close()
		logger.Infof("receiver closed, error=%v", err)

		os.Exit(0)
	}()

	for {
		runtime.Gosched()
	}
}
