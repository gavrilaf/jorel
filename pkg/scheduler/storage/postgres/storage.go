package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/gavrilaf/dyson/pkg/dlog"
	base "github.com/gavrilaf/dyson/pkg/scheduler/storage"
)

type storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, databaseUrl string) (base.SchedulerStorage, error) {
	log := logrusadapter.NewLogger(dlog.FromContext(ctx))

	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url, %w", err)
	}

	poolConfig.ConnConfig.Logger = log

	db, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("faild to create connection pool, %w", err)
	}

	return &storage{
		db: db,
	}, nil
}


func (s *storage) Save(ctx context.Context, scheduledTime time.Time, msg base.Message) error {
	attributes, err := json.Marshal(msg.Attributes)
	if err != nil {
		return fmt.Errorf("unable to marshal attributes, %w", msg.Attributes)
	}

	_, err = s.db.Exec(ctx, "insert into messages(scheduledtime, done, data, attributes) values($1, $2, $3, $4)",
		scheduledTime, false, msg.Data, attributes)
	return err
}

func (s *storage) GetMessages(ctx context.Context, olderThan time.Time, limit int) ([]base.ScheduledMessage, error) {
	return nil, nil
}

func (s *storage) Close() error {
	s.db.Close()
	return nil
}
