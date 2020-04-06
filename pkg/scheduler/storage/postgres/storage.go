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

func (s *storage) GetLatest(ctx context.Context, olderThan time.Time, handler base.Handler) (bool, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to start transaction, %w", err)
	}

	sql := "UPDATE messages SET done = 't' WHERE id = (SELECT id " +
				"FROM messages WHERE done = 'f' AND scheduledtime <= $1 " +
				"ORDER BY scheduledtime FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING scheduledtime, data, attributes;"

	rows, err := tx.Query(ctx, sql, olderThan)
	if err != nil {
		tx.Rollback(ctx)
		return false, fmt.Errorf("failed to query scheduled items, %w", err)
	}

	handled := false
	if rows.Next() {
		var scheduledTime time.Time
		var data []byte
		var attrs []byte

		err = rows.Scan(&scheduledTime, &data, &attrs)
		if err != nil {
			tx.Rollback(ctx)
			return false, fmt.Errorf("failed to query scheduled items, %w", err)
		}

		var parsedAttrs map[string]string
		err = json.Unmarshal(attrs, &parsedAttrs)
		if err != nil {
			tx.Rollback(ctx)
			return false, fmt.Errorf("failed to unmarshal attributes, %w", err)
		}

		msg := base.ScheduledMessage{
			Message:       base.Message{
				Data:       data,
				Attributes: parsedAttrs,
			},
			ScheduledTime: scheduledTime.UTC(),
		}

		err := handler.Handle(ctx, msg)
		if err != nil {
			tx.Rollback(ctx)
			return false, fmt.Errorf("failed to handler message, %w", err)
		}

		handled = true
	}

	rows.Close()

	err = tx.Commit(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to commit transaction, %w", err)
	}

	return handled, nil
}

func (s *storage) Close() error {
	s.db.Close()
	return nil
}