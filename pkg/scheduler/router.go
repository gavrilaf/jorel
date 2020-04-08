package scheduler

import (
	"context"
	"github.com/gavrilaf/dyson/pkg/dlog"
	"io"

	"github.com/gavrilaf/dyson/pkg/msgqueue"
)

//go:generate mockery -name Router -outpkg schedulermocks -output ./schedulermocks -dir .
type Router interface {
	io.Closer
	Publish(ctx context.Context, msgType string, data []byte, attributes map[string]string) (msgqueue.PublishResult, error)
}

type routerImpl struct {
	defaultEgress msgqueue.Publisher
	routes        map[string]msgqueue.Publisher
}

func NewRouter(ctx context.Context, config Config) (Router, error) {
	defaultEgress, err := msgqueue.NewPublisher(ctx, config.ProjectID, config.DefaultEgress.Name)
	if err != nil {
		return &routerImpl{}, err
	}

	routes := make(map[string]msgqueue.Publisher)
	for m, r := range config.Routing {
		p, err := msgqueue.NewPublisher(ctx, config.ProjectID, r.Name)
		if err != nil {
			return &routerImpl{}, err
		}

		routes[m] = p
	}

	return &routerImpl{
		defaultEgress: defaultEgress,
		routes:        routes,
	}, nil
}

func (r *routerImpl) Publish(ctx context.Context, msgType string, data []byte, attributes map[string]string) (msgqueue.PublishResult, error) {
	p, ok := r.routes[msgType]
	if !ok {
		p = r.defaultEgress
	}

	res, err := p.Publish(ctx, data, attributes)
	if err != nil {
		dlog.FromContext(ctx).WithError(err).Errorf("failed to publish message to the topic %s", p.TopicID())
	} else {
		dlog.FromContext(ctx).Infof("message published to the topic %s", p.TopicID())
	}

	return res, err
}

func (r *routerImpl) Close() error {
	err := r.defaultEgress.Close()
	for _, e := range r.routes {
		if err2 := e.Close(); err2 != nil {
			err = err2
		}
	}

	return err
}
