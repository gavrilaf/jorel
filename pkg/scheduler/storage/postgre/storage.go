package postgre

import (
	"context"
	"time"

	"github.com/gavrilaf/dyson/pkg/scheduler/storage"
)

type Storage struct{}

func (s *Storage) Save(ctx context.Context, scheduledTime time.Time, msg storage.Message) error {
	return nil
}

func (s *Storage) GetMessages(ctx context.Context, olderThan time.Time, limit int) ([]storage.ScheduledMessage, error) {
	return nil, nil
}
