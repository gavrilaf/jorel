package storage

import (
	"context"
	"io"
	"time"
)

type Message struct {
	Data []byte
	Attributes map[string]string
}

type ScheduledMessage struct {
	Message
	ScheduledTime time.Time
}

//go:generate mockery -name Handler -outpkg storagemocks -output ./storagemocks -dir .
type Handler interface {
	Handle(ctx context.Context, msg ScheduledMessage) error
}

//go:generate mockery -name SchedulerStorage -outpkg storagemocks -output ./storagemocks -dir .
type SchedulerStorage interface {
	io.Closer
	Save(ctx context.Context, scheduledTime time.Time, msg Message) error
	GetLatest(ctx context.Context, olderThan time.Time, handler Handler) (bool, error)
}
