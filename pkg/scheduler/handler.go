package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage"
)

type HandlerConfig struct {
	Publisher  msgqueue.Publisher
	Storage    storage.SchedulerStorage
	TimeSource TimeSource
}

type Handler struct {
	publisher  msgqueue.Publisher
	storage    storage.SchedulerStorage
	timeSource TimeSource
}

func NewHandler(config HandlerConfig) *Handler {
	return &Handler{
		publisher:  config.Publisher,
		storage:    config.Storage,
		timeSource: config.TimeSource,
	}
}

func (h *Handler) Receive(ctx context.Context, data []byte, attributes map[string]string) error {
	msgAttributes, err := msgqueue.NewMsgAttributes(attributes)
	if err != nil {
		return fmt.Errorf("failed to parse attributes, %w", err)
	}

	if msgAttributes.DelayInSeconds == 0 {
		return h.publish(ctx, data, msgAttributes.Original)
	} else {
		return h.save(ctx, data, msgAttributes)
	}
}

func (h *Handler) publish(ctx context.Context, data []byte, attributes map[string]string) error {
	result, err := h.publisher.Publish(ctx, data, attributes)
	if err != nil {
		return fmt.Errorf("failed to publish message, %w", err)
	}

	msgID, err := result.GetMessageID(ctx)
	if err != nil {
		return fmt.Errorf("failed to read message id, %w", err)
	}

	dlog.FromContext(ctx).Infof("Resend message with id: %s", msgID)

	return nil
}

func (h *Handler) save(ctx context.Context, data []byte, msgAttributes msgqueue.MsgAttributes) error {
	scheduledTime := h.timeSource.Now().Add(msgAttributes.DelayInSeconds * time.Second)

	message := storage.Message{
		Data:       data,
		Attributes: msgAttributes.Original,
	}

	return h.storage.Save(ctx, scheduledTime, message)
}
