package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/gavrilaf/dyson/pkg/dlog"
	"github.com/gavrilaf/dyson/pkg/msgqueue"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage"
)

type IngressConfig struct {
	Router     Router
	Storage    storage.SchedulerStorage
	TimeSource TimeSource
}

type Ingress struct {
	router     Router
	storage    storage.SchedulerStorage
	timeSource TimeSource
}

func NewIngress(config IngressConfig) *Ingress {
	return &Ingress{
		router:     config.Router,
		storage:    config.Storage,
		timeSource: config.TimeSource,
	}
}

func (p *Ingress) Receive(ctx context.Context, data []byte, attributes map[string]string) error {
	msgAttributes, err := msgqueue.NewMsgAttributes(attributes)
	if err != nil {
		return fmt.Errorf("failed to parse attributes, %w", err)
	}

	dlog.FromContext(ctx).Infof("Received message with delay: %d", msgAttributes.DelayInSeconds)

	if msgAttributes.DelayInSeconds == 0 {
		_, err = p.router.Publish(ctx, msgAttributes.MessageType, data, msgAttributes.Original)
		if err != nil {
			return fmt.Errorf("failed to resend message, %w", err)
		}
	} else {
		return p.save(ctx, data, msgAttributes)
	}

	return nil
}

func (p *Ingress) save(ctx context.Context, data []byte, msgAttributes msgqueue.MsgAttributes) error {
	scheduledTime := p.timeSource.Now().Add(time.Duration(msgAttributes.DelayInSeconds) * time.Second)

	message := storage.Message{
		Data:       data,
		Attributes: msgAttributes.Original,
	}

	dlog.FromContext(ctx).Infof("Save message with scheduled time: %v, messge type: %s", scheduledTime, msgAttributes.MessageType)
	return p.storage.Save(ctx, scheduledTime, msgAttributes.MessageType, msgAttributes.AggregationID, message)
}
