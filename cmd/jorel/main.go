package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/scheduler"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage/postgres"
	"github.com/gavrilaf/dyson/pkg/utils"
)

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		logger.Panicf("unable to read config file, %v", err)
	}

	config, err := scheduler.ParseConfig(content)
	if err != nil {
		logger.Panicf("unable to parse config file, %v", err)
	}

	// db
	dbUrl := os.Getenv("JOR_EL_POSTGRES_URL")
	storage, err := postgres.NewStorage(ctx, dbUrl)
	if err != nil {
		logger.Panicf("failed to connect database, %v", err)
	}

	// router
	router, err := scheduler.NewRouter(ctx, config)
	if err != nil {
		logger.Panicf("failed to create router, %v", err)
	}

	// ingress
	logger.Info("Starting ingress")

	ingressReceiver, err := msgqueue.NewReceiver(ctx, config.ProjectID, config.IngressSubscription)
	if err != nil {
		logger.Panicf("failed to create ingress receiver, %v", err)
	}

	ingressConfig := scheduler.IngressConfig{
		Router:  router,
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
		Router:  router,
		Storage:    storage,
		TimeSource: scheduler.SystemTime{},
	}

	ticker := scheduler.NewTicker(tickerConfig)

	cancelCtx, cancelFn := context.WithCancel(ctx)
	ticker.RunTicker(cancelCtx)

	// done

	logger.Info("Jor-el started")

	utils.WaitForShutdown(ctx, func() {
		cancelFn()

		err := ingressReceiver.Close()
		logger.Infof("ingressReceiver closed, error=%v", err)

		err = router.Close()
		logger.Infof("publisher closed, error=%v", err)

		err = storage.Close()
		logger.Infof("storage closed, error=%v", err)
	})
}
