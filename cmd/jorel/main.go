package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/scheduler"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/testdata"
)

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	receiver, err := msgqueue.NewReceiver(ctx, testdata.ProjectID, testdata.IngressSubs)
	if err != nil {
		logger.Panicf("failed to create receiver, %v", err)
	}

	publisher, err := msgqueue.NewPublisher(ctx, testdata.ProjectID, testdata.EgressTopic)
	if err != nil {
		logger.Panicf("failed to create publisher, %v", err)
	}

	logger.Info("Starting jor-el")

	config := scheduler.HandlerConfig{
		Publisher: publisher,
		TimeSource: scheduler.SystemTime{},
	}

	handler := scheduler.NewHandler(config)

	err = receiver.Run(ctx, handler)
	if err != nil {
		logger.Panicf("failed to run receiver, %v", err)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info("Ctrl+C pressed in Terminal")

		err := receiver.Close()
		logger.Infof("receiver closed, error=%v", err)

		err = publisher.Close()
		logger.Infof("publisher closed, error=%v", err)

		os.Exit(0)
	}()

	for {
		time.Sleep(time.Second)
	}
}
