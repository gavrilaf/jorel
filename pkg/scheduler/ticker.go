package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage"
)

const (
	tickPeriod    = 1 * time.Second
	deviationTime = 2 * time.Second
)

type TickerConfig struct {
	Router     Router
	Storage    storage.SchedulerStorage
	TimeSource TimeSource
}

type Ticker struct {
	router     Router
	storage    storage.SchedulerStorage
	timeSource TimeSource
}

func NewTicker(config TickerConfig) *Ticker {
	return &Ticker{
		router:     config.Router,
		storage:    config.Storage,
		timeSource: config.TimeSource,
	}
}

// Activates ticker loop
func (h *Ticker) RunTicker(ctx context.Context) {
	ticker := time.NewTicker(tickPeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				dlog.FromContext(ctx).Info("ticker is shut down")
				return
			case <-ticker.C:
				h.handleTick(ctx)
			}
		}

	}()
}

func (t *Ticker) handleTick(ctx context.Context) {
	scanTime := t.timeSource.Now().Add(deviationTime)

	counter := 0
	continueHandling := true
	var err error

	for continueHandling {
		continueHandling, err = t.storage.GetLatest(ctx, scanTime, t)
		if err != nil {
			dlog.FromContext(ctx).WithError(err).Error("failed to handle message")
			break
		}
		if continueHandling {
			counter += 1
		}
	}

	dlog.FromContext(ctx).Infof("handled %d messages, scan time: %v", counter, scanTime)
}

func (t *Ticker) HandleMessage(ctx context.Context, msg storage.ScheduledMessage) error {
	_, err := t.router.Publish(ctx, msg.MessageType, msg.Data, msg.Attributes)
	if err != nil {
		return fmt.Errorf("failed to pubslish message, %w", err)
	}

	dlog.FromContext(ctx).Infof("messaged published, current time: %v, scheduled time: %v", t.timeSource.Now(), msg.ScheduledTime)

	return nil
}
