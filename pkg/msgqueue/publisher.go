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
	TopicID() string
	Publish(ctx context.Context, data []byte, attributes map[string]string) (PublishResult, error)
}

func NewPublisher(ctx context.Context, projectID string, topicID string) (Publisher, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic := client.Topic(topicID)

	return &publisherImpl{
		client:  client,
		topic:   topic,
		topicID: topicID,
	}, nil
}

// publisherImpl

type publisherImpl struct {
	client  *pubsub.Client
	topic   *pubsub.Topic
	topicID string
}

func (p *publisherImpl) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}

func (p *publisherImpl) TopicID() string {
	return p.topicID
}

func (p *publisherImpl) Publish(ctx context.Context, data []byte, attributes map[string]string) (PublishResult, error) {
	msg := &pubsub.Message{Data: data, Attributes: attributes}
	return &publishResult{result: p.topic.Publish(ctx, msg)}, nil
}
