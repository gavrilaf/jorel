package scheduler

import (
	"context"
	"fmt"
	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
)

type HandlerConfig struct {
	Receiver msgqueue.Receiver
	Publisher msgqueue.Publisher
}


type Handler struct {
	receiver msgqueue.Receiver
	publisher msgqueue.Publisher
}

func NewHandler(config HandlerConfig) *Handler {
	return &Handler{
		receiver:  config.Receiver,
		publisher: config.Publisher,
	}
}

func (h *Handler) Run(ctx context.Context) error {
	return h.receiver.Run(ctx, h.handle)
}

func (h *Handler) handle(ctx context.Context, data []byte, attributes map[string]string) error {
	msgAttributes, err := msgqueue.NewMsgAttributes(attributes)
	if err != nil {
		return fmt.Errorf("failed to parse attributes, %w", err)
	}

	result, err := h.publisher.Publish(ctx, data, msgAttributes)
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