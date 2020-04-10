package postgres_test

import (
	"context"
	"encoding/json"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gavrilaf/jorel/pkg/scheduler/storage"
	"github.com/gavrilaf/jorel/pkg/scheduler/storage/postgres"
	"github.com/gavrilaf/jorel/pkg/scheduler/storage/storagemocks"
)

func TestStorageBase(t *testing.T) {
	ctx := context.Background()
	dbUrl := os.Getenv("JOR_EL_POSTGRES_TEST_URL")

	stg, err := postgres.NewStorage(ctx, dbUrl)
	assert.NoError(t, err)

	now := time.Now().UTC()

	fiveSecs := now.Add(5 * time.Second)
	tenSecs := now.Add(10 * time.Second)

	msgWithAttributes := storage.Message{
		Data:       []byte("123"),
		Attributes: map[string]string{"one": "two"},
	}

	msgWithEmptyAttributes := storage.Message{
		Data: []byte("1234567"),
	}

	err = stg.Save(ctx, fiveSecs, "cancel", "", msgWithAttributes)
	assert.NoError(t, err)

	err = stg.Save(ctx, tenSecs, "", "", msgWithEmptyAttributes)
	assert.NoError(t, err)

	handler := &storagemocks.Handler{}
	handler.On("HandleMessage", mock.Anything, mock.Anything).Return(nil)

	scanTime := now.Add(15 * time.Second)

	handled, err := stg.GetLatest(ctx, scanTime, handler)
	assert.NoError(t, err)
	assert.True(t, handled)

	expectedMessage := storage.ScheduledMessage{
		Message:       msgWithAttributes,
		MessageType:   "cancel",
		ScheduledTime: fiveSecs,
	}
	handler.AssertCalled(t, "HandleMessage", mock.Anything, expectedMessage)

	handled, err = stg.GetLatest(ctx, scanTime, handler)
	assert.NoError(t, err)
	assert.True(t, handled)

	expectedMessage = storage.ScheduledMessage{
		Message:       msgWithEmptyAttributes,
		ScheduledTime: tenSecs,
	}
	handler.AssertCalled(t, "HandleMessage", mock.Anything, expectedMessage)

	handled, err = stg.GetLatest(ctx, scanTime, handler)
	assert.NoError(t, err)
	assert.False(t, handled)

	handler.AssertNumberOfCalls(t, "HandleMessage", 2)

	t.Logf("cleaning up...")
	err = cleanupDb(ctx)
	assert.NoError(t, err)
}

func TestStorageBigObject(t *testing.T) {
	ctx := context.Background()
	dbUrl := os.Getenv("JOR_EL_POSTGRES_TEST_URL")

	stg, err := postgres.NewStorage(ctx, dbUrl)
	assert.NoError(t, err)

	now := time.Now().UTC()
	scheduledTime := now.Add(24 * time.Hour)

	bigData := make([]int, 2048)
	for i := range bigData {
		bigData[i] = i
	}

	object := struct {
		ID          string
		SomeBigData []int
		SomeMap     map[string]interface{}
	}{
		ID:          "12345",
		SomeBigData: bigData,
		SomeMap:     map[string]interface{}{"one": "two", "three": 3, "four": bigData},
	}

	marshaledObject, err := json.Marshal(&object)
	assert.NoError(t, err)

	message := storage.Message{
		Data:       marshaledObject,
		Attributes: map[string]string{"one": "two", "2": "3", "message-id": "12345"},
	}

	err = stg.Save(ctx, scheduledTime, "big-object", "", message)
	assert.NoError(t, err)

	handler := &storagemocks.Handler{}
	handler.On("HandleMessage", mock.Anything, mock.Anything).Return(nil)

	scanTime := scheduledTime.Add(1 * time.Second)

	handled, err := stg.GetLatest(ctx, scanTime, handler)
	assert.NoError(t, err)
	assert.True(t, handled)

	expectedMessage := storage.ScheduledMessage{
		Message:       message,
		MessageType:   "big-object",
		ScheduledTime: scheduledTime,
	}
	handler.AssertCalled(t, "HandleMessage", mock.Anything, expectedMessage)

	handled, err = stg.GetLatest(ctx, scanTime, handler)
	assert.NoError(t, err)
	assert.False(t, handled)

	handler.AssertNumberOfCalls(t, "HandleMessage", 1)

	t.Logf("cleaning up...")
	err = cleanupDb(ctx)
	assert.NoError(t, err)
}

func TestStorageConcurrency(t *testing.T) {
	ctx := context.Background()
	dbUrl := os.Getenv("JOR_EL_POSTGRES_TEST_URL")

	stg, err := postgres.NewStorage(ctx, dbUrl)
	assert.NoError(t, err)

	now := time.Now().UTC()
	fiveSecs := now.Add(5 * time.Second)
	scanTime := now.Add(15 * time.Second)

	msg := storage.Message{
		Data:       []byte("123"),
		Attributes: map[string]string{"one": "two"},
	}

	objectsCount := 20

	for i := 0; i < objectsCount; i++ {
		err = stg.Save(ctx, fiveSecs, "update", "", msg)
		assert.NoError(t, err)
	}

	var wg sync.WaitGroup
	counter := 0

	handleFn := func() {
		continueHandling := true
		var err error

		for continueHandling {
			handler := &storagemocks.Handler{}
			handler.On("HandleMessage", mock.Anything, mock.Anything).Return(nil)

			continueHandling, err = stg.GetLatest(ctx, scanTime, handler)
			assert.NoError(t, err)

			if continueHandling {
				counter += 1
			}

			runtime.Gosched()
		}

		wg.Done()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go handleFn()
	}

	wg.Wait()

	assert.Equal(t, objectsCount, counter)

	handler := &storagemocks.Handler{}
	handler.On("HandleMessage", mock.Anything, mock.Anything).Return(nil)

	handled, err := stg.GetLatest(ctx, scanTime, handler)
	assert.NoError(t, err)
	assert.False(t, handled)
	handler.AssertNumberOfCalls(t, "HandleMessage", 0)

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
