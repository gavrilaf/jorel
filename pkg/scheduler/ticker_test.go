package scheduler_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gavrilaf/dyson/pkg/msgqueue/msgqueuemocks"
	"github.com/gavrilaf/dyson/pkg/scheduler"
	"github.com/gavrilaf/dyson/pkg/scheduler/schedulermocks"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage/storagemocks"
)

func TestTickerReceive(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("storage does not contain eligible messages", func(t *testing.T) {
		ticker, mocks := tickerWithMocks()

		mocks.timeSource.On("Now").Return(now)
		mocks.storage.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

		err := ticker.Tick(ctx)
		assert.NoError(t, err)

		mocks.storage.AssertCalled(t, "GetLatest", mock.Anything, now.Add(1 * time.Second), ticker)
	})

	t.Run("storage contains eligible messages", func(t *testing.T) {
		ticker, mocks := tickerWithMocks()

		counter := 0

		mocks.timeSource.On("Now").Return(now)
		mocks.storage.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(func(context.Context, time.Time, storage.Handler) bool {
			if counter >= 2 {
				return false
			}
			counter += 1
			return true
		}, nil)

		err := ticker.Tick(ctx)
		assert.NoError(t, err)

		mocks.storage.AssertCalled(t, "GetLatest", mock.Anything, now.Add(1 * time.Second), ticker)
		mocks.storage.AssertNumberOfCalls(t, "GetLatest", 3)
	})

	t.Run("storage fails", func(t *testing.T) {
		ticker, mocks := tickerWithMocks()


		mocks.timeSource.On("Now").Return(now)
		mocks.storage.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(false, fmt.Errorf(""))

		err := ticker.Tick(ctx)
		assert.Error(t, err)
	})
}

func TestTickerHandleMessage(t *testing.T) {
	ctx := context.Background()

	publishResult := &msgqueuemocks.PublishResult{}

	now := time.Now()

	data := []byte("123")
	attributes := map[string]string{"one": "two"}

	msg := storage.ScheduledMessage{
		Message:       storage.Message{
			Data: data,
			Attributes: attributes,
		},
	}

	t.Run("publish message", func(t *testing.T) {
		ticker, mocks := tickerWithMocks()

		mocks.timeSource.On("Now").Return(now)
		mocks.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := ticker.HandleMessage(ctx, msg)
		assert.NoError(t, err)

		mocks.publisher.AssertCalled(t, "Publish", mock.Anything, data, attributes)
	})
}

type tickerMocks struct {
	publisher  *msgqueuemocks.Publisher
	storage    *storagemocks.SchedulerStorage
	timeSource *schedulermocks.TimeSource
}

func tickerWithMocks() (*scheduler.Ticker, tickerMocks) {
	m := tickerMocks{
		publisher:  &msgqueuemocks.Publisher{},
		storage:    &storagemocks.SchedulerStorage{},
		timeSource: &schedulermocks.TimeSource{},
	}

	ticker := scheduler.NewTicker(scheduler.TickerConfig{
		Publisher:  m.publisher,
		Storage:    m.storage,
		TimeSource: m.timeSource,
	})

	return ticker, m
}
