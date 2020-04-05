package msgqueue

import (
	"context"
	"io"

	"cloud.google.com/go/pubsub"
)

// PublishResult

//go:generate mockery -name PublishResult -outpkg msgqueuemocks -output ./msgqueuemocks -dir .
type PublishResult interface {
	GetMessageID(ctx context.Context) (string, error)
}

type publishResult struct {
	result *pubsub.PublishResult
}

func (p *publishResult) GetMessageID(ctx context.Context) (string, error) {
	return p.result.Get(ctx)
}

// Publisher

//go:generate mockery -name Publisher -outpkg msgqueuemocks -output ./msgqueuemocks -dir .
type Publisher interface {
	io.Closer
	Publish(ctx context.Context, data []byte, attributes map[string]string) (PublishResult, error)
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

func (c *publisherImpl) Publish(ctx context.Context, data []byte, attributes map[string]string) (PublishResult, error) {
	msg := &pubsub.Message{Data: data, Attributes: attributes}
	return &publishResult{result: c.topic.Publish(ctx, msg)}, nil
}
