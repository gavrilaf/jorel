package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gavrilaf/jorel/pkg/dlog"
)

func WaitForShutdown(ctx context.Context, cleanFn func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		dlog.FromContext(ctx).Info("Ctrl+C pressed in Terminal")

		cleanFn()

		os.Exit(0)
	}()

	for {
		time.Sleep(time.Second)
	}
}
