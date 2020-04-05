package storage

import "time"

type Message struct {
	Data []byte
	Attributes map[string]string
}

type ScheduledMessage struct {
	Message
	ScheduledTime time.Time
}

//go:generate mockery -name SchedulerStorage -outpkg storagemocks -output ./storagemocks -dir .
type SchedulerStorage interface {
	Save(scheduledTime time.Time, msg Message) error
	GetMessages(olderThan time.Time, limit int) ([]ScheduledMessage, error)
}
