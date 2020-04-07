package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage"
)

const deviationTimeInSeconds = 2

type TickerConfig struct {
	Publisher  msgqueue.Publisher
	Storage    storage.SchedulerStorage
	TimeSource TimeSource
}

type Ticker struct {
	publisher  msgqueue.Publisher
	storage    storage.SchedulerStorage
	timeSource TimeSource
}

func NewTicker(config TickerConfig) *Ticker {
	return &Ticker{
		publisher:  config.Publisher,
		storage:    config.Storage,
		timeSource: config.TimeSource,
	}
}

// Activates ticker loop
func (h *Ticker) Tick(ctx context.Context) error {
	return h.handleTick(ctx)
}

func (h *Ticker) handleTick(ctx context.Context) error {
	scanTime := h.timeSource.Now().Add(deviationTimeInSeconds * time.Second)

	counter := 0
	continueHandling := true
	var err error

	for continueHandling {
		continueHandling, err = h.storage.GetLatest(ctx, scanTime, h)
		if err != nil {
			dlog.FromContext(ctx).WithError(err).Error("failed to handle message")
			break
		}
		if continueHandling {
			counter += 1
		}
	}

	dlog.FromContext(ctx).Infof("handled %d messages, scan time: %v", counter, scanTime)
	return err
}

func (h *Ticker) HandleMessage(ctx context.Context, msg storage.ScheduledMessage) error {
	_, err := h.publisher.Publish(ctx, msg.Data, msg.Attributes)
	if err != nil {
		return fmt.Errorf("failed to pubslish message, %w", err)
	}

	dlog.FromContext(ctx).Infof("messaged published, current time: %v, scheduled time: %v", h.timeSource.Now(), msg.ScheduledTime)

	return nil
}
