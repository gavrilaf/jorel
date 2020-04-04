package main

import (
	"context"

	"github.com/gavrilaf/dyson/pkg/dlog"
)

func main() {
	ctx := context.Background()
	logger := dlog.FromContext(ctx)

	logger.Info("Starting capacitor")
}
