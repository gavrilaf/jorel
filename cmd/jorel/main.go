package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/scheduler"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage/postgres"
	"github.com/gavrilaf/dyson/pkg/testdata"
)

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	publisher, err := msgqueue.NewPublisher(ctx, testdata.ProjectID, testdata.EgressTopic)
	if err != nil {
		logger.Panicf("failed to create publisher, %v", err)
	}

	dbUrl := os.Getenv("JOR_EL_POSTGRE_URL")
	storage, err := postgres.NewStorage(ctx, dbUrl)
	if err != nil {
		logger.Panicf("failed to connect database, %v", err)
	}

	// ingress

	logger.Info("Starting ingress")

	ingressReceiver, err := msgqueue.NewReceiver(ctx, testdata.ProjectID, testdata.IngressSubs)
	if err != nil {
		logger.Panicf("failed to create ingressReceiver, %v", err)
	}

	ingressConfig := scheduler.IngressConfig{
		Publisher:  publisher,
		Storage:    storage,
		TimeSource: scheduler.SystemTime{},
	}

	ingress := scheduler.NewIngress(ingressConfig)

	err = ingressReceiver.Run(ctx, ingress)
	if err != nil {
		logger.Panicf("failed to run ingress receiver, %v", err)
	}

	// ticker

	logger.Info("Starting ticker")

	tickerConfig := scheduler.TickerConfig{
		Publisher:  publisher,
		Storage:    storage,
		TimeSource: scheduler.SystemTime{},
	}

	ticker := scheduler.NewTicker(tickerConfig)

	done := false

	go func() {
		for !done {
			ticker.Tick(ctx)
			time.Sleep(time.Second)
		}
	}()

	// done

	logger.Info("Jor-el started")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info("Ctrl+C pressed in Terminal")

		done = true

		err := ingressReceiver.Close()
		logger.Infof("ingressReceiver closed, error=%v", err)

		err = publisher.Close()
		logger.Infof("publisher closed, error=%v", err)

		err = storage.Close()
		logger.Infof("storage closed, error=%v", err)

		os.Exit(0)
	}()

	for {
		time.Sleep(time.Second)
	}
}
