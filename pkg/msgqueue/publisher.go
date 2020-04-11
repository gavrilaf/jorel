package msgqueue

import (
	"context"
	"io"

	"cloud.google.com/go/pubsub"
)

// PublisherFactory

//go:generate mockery -name PublisherFactory -outpkg msgqueuemocks -output ./msgqueuemocks -dir .
type PublisherFactory interface {
	NewPublisher(ctx context.Context, projectID string, topicID string) (Publisher, error)
}

func NewPublisherFactory() PublisherFactory {
	return &publisherFactory{}
}

type publisherFactory struct {}

func (publisherFactory) NewPublisher(ctx context.Context, projectID string, topicID string) (Publisher, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic := client.Topic(topicID)

	return &publisher{
		client:  client,
		topic:   topic,
		topicID: topicID,
	}, nil
}

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

type publisher struct {
	client  *pubsub.Client
	topic   *pubsub.Topic
	topicID string
}

func (p *publisher) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}

func (p *publisher) TopicID() string {
	return p.topicID
}

func (p *publisher) Publish(ctx context.Context, data []byte, attributes map[string]string) (PublishResult, error) {
	msg := &pubsub.Message{Data: data, Attributes: attributes}
	return &publishResult{result: p.topic.Publish(ctx, msg)}, nil
}
