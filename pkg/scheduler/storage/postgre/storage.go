package postgre

import (
	"context"
	"fmt"
	"time"

	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	base "github.com/gavrilaf/dyson/pkg/scheduler/storage"
	"github.com/gavrilaf/dyson/pkg/dlog"
)

type storage struct{
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, databaseUrl string) (base.SchedulerStorage, error) {

	log := logrusadapter.NewLogger(dlog.FromContext(ctx))

	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url, %w", err)
	}

	poolConfig.ConnConfig.Logger = log

	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("faild to create connection pool, %w", err)
	}

	return &storage{
		db: db,
	}, nil
}

/*
	Table(
		id
		scheduledTime
		data
		attributes
	)

 */

func (s *storage) Save(ctx context.Context, scheduledTime time.Time, msg base .Message) error {
	return nil
}

func (s *storage) GetMessages(ctx context.Context, olderThan time.Time, limit int) ([]base.ScheduledMessage, error) {
	return nil, nil
}

func (s *storage) Close() error {
	s.db.Close()
	return nil
}
