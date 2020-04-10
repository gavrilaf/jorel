package msgqueue

import (
	"context"
	"io"

	"cloud.google.com/go/pubsub"

	"github.com/gavrilaf/jorel/pkg/dlog"
)

type Handler interface {
	Receive(ctx context.Context, data []byte, attributes map[string]string) error
}

type Receiver interface {
	io.Closer
	Run(ctx context.Context, handler Handler) error
}

func NewReceiver(ctx context.Context, projectID string, subscriptionID string) (Receiver, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	sub := client.Subscription(subscriptionID)

	return &receiverImpl{
		client: client,
		sub:    sub,
	}, nil
}

// receiverImpl

type receiverImpl struct {
	client *pubsub.Client
	sub    *pubsub.Subscription
}

func (p *receiverImpl) Run(ctx context.Context, handler Handler) error {
	go func() {
		err := p.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			defer msg.Ack()

			log := dlog.FromContext(ctx).WithFields(map[string]interface{}{"message-id": msg.ID, "attributes": msg.Attributes})
			log.Info("received message")

			ctxWithLog := dlog.WithLogger(ctx, log)
			err := handler.Receive(ctxWithLog, msg.Data, msg.Attributes)

			if err != nil {
				log.WithError(err).Error("failed to handle message")
			} else {
				log.Info("message handled")
			}
		})

		if err != nil {
			dlog.FromContext(ctx).WithError(err).Panic("pub/sub subscription Receive error")
		}
	}()

	return nil
}

func (p *receiverImpl) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}
