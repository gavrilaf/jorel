package postgres_test

import (
	"context"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage"
	"github.com/jackc/pgx/v4"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/dyson/pkg/scheduler/storage/postgres"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	dbUrl := os.Getenv("JOR_EL_POSTGRES_TEST_URL")

	stg, err := postgres.NewStorage(ctx, dbUrl)
	assert.NoError(t, err)

	now := time.Now()

	fiveSecs := now.Add(5 * time.Second)
	tenSecs := now.Add(10 * time.Second)

	msgWithAttributes := storage.Message{
		Data:       []byte("123"),
		Attributes: map[string]string{"one": "two"},
	}

	msgWithEmptyAttributes := storage.Message{
		Data:       []byte("1234567"),
	}

	err = stg.Save(ctx, fiveSecs, msgWithAttributes)
	assert.NoError(t, err)

	err = stg.Save(ctx, tenSecs, msgWithEmptyAttributes)
	assert.NoError(t, err)

	t.Logf("cleaning up...")
	err = cleanupDb(ctx)
	assert.NoError(t, err)
}

func cleanupDb(ctx context.Context) error {
	dbUrl := os.Getenv("JOR_EL_POSTGRES_TEST_URL")

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, "delete from messages")
	return err
}

