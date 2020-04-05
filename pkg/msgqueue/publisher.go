package msgqueue

import (
	"context"
	"io"

	"cloud.google.com/go/pubsub"
)

// PublishResult

type PublishResult struct {
	result *pubsub.PublishResult
}

func (p *PublishResult) GetMessageID(ctx context.Context) (string, error) {
	return p.result.Get(ctx)
}

// Publisher

type Publisher interface {
	io.Closer
	Publish(ctx context.Context, data []byte, attributes MsgAttributes) (*PublishResult, error)
}

func NewPublisher(ctx context.Context, projectID string, topicID string) (Publisher, error) {
	pubsub, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic := pubsub.Topic(topicID)

	return &publisherImpl{
		pubsub: pubsub,
		topic:  topic,
	}, nil
}

// publisherImpl

type publisherImpl struct {
	pubsub *pubsub.Client
	topic  *pubsub.Topic
}

func (c *publisherImpl) Close() error {
	if c.pubsub != nil {
		return c.pubsub.Close()
	}
	return nil
}

func (c *publisherImpl) Publish(ctx context.Context, data []byte, attributes MsgAttributes) (*PublishResult, error) {
	msg := &pubsub.Message{Data: data, Attributes: attributes.GetAttributes()}
	return &PublishResult{result: c.topic.Publish(ctx, msg)}, nil
}
