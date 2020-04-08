package scheduler_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
		cancelCtx, cancel := context.WithCancel(ctx)
		ticker, mocks := tickerWithMocks()
		hook := test.NewGlobal()

		mocks.timeSource.On("Now").Return(now)
		mocks.storage.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

		ticker.RunTicker(cancelCtx)

		time.Sleep(1 * time.Second)
		cancel()
		time.Sleep(1 * time.Second)

		mocks.storage.AssertCalled(t, "GetLatest", mock.Anything, now.Add(2*time.Second), ticker)

		expectedLogs := []log{
			{
				msg:   fmt.Sprintf("handled 0 messages, scan time: %v", now.Add(2*time.Second)),
				level: logrus.InfoLevel,
			},
			{
				msg:   "ticker is shut down",
				level: logrus.InfoLevel,
			},
		}
		assertLogs(t, expectedLogs, hook.AllEntries())
		hook.Reset()

	})

	t.Run("storage contains eligible messages", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		ticker, mocks := tickerWithMocks()
		hook := test.NewGlobal()

		counter := 0

		mocks.timeSource.On("Now").Return(now)
		mocks.storage.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(func(context.Context, time.Time, storage.Handler) bool {
			if counter >= 2 {
				return false
			}
			counter += 1
			return true
		}, nil)

		ticker.RunTicker(cancelCtx)

		time.Sleep(1 * time.Second)
		cancel()
		time.Sleep(1 * time.Second)

		mocks.storage.AssertCalled(t, "GetLatest", mock.Anything, now.Add(2*time.Second), ticker)
		mocks.storage.AssertNumberOfCalls(t, "GetLatest", 3)

		expectedLogs := []log{
			{
				msg:   fmt.Sprintf("handled 2 messages, scan time: %v", now.Add(2*time.Second)),
				level: logrus.InfoLevel,
			},
			{
				msg:   "ticker is shut down",
				level: logrus.InfoLevel,
			},
		}
		assertLogs(t, expectedLogs, hook.AllEntries())
		hook.Reset()
	})

	t.Run("storage fails", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		ticker, mocks := tickerWithMocks()
		hook := test.NewGlobal()

		mocks.timeSource.On("Now").Return(now)
		mocks.storage.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(false, fmt.Errorf("error"))

		ticker.RunTicker(cancelCtx)

		time.Sleep(1 * time.Second)
		cancel()
		time.Sleep(1 * time.Second)

		expectedLogs := []log{
			{
				msg:   "failed to handle message",
				level: logrus.ErrorLevel,
			},
			{
				msg:   fmt.Sprintf("handled 0 messages, scan time: %v", now.Add(2*time.Second)),
				level: logrus.InfoLevel,
			},
			{
				msg:   "ticker is shut down",
				level: logrus.InfoLevel,
			},
		}
		assertLogs(t, expectedLogs, hook.AllEntries())
		hook.Reset()
	})
}

func TestTickerHandleMessage(t *testing.T) {
	ctx := context.Background()

	publishResult := &msgqueuemocks.PublishResult{}

	now := time.Now()

	data := []byte("123")
	attributes := map[string]string{"one": "two"}

	msg := storage.ScheduledMessage{
		Message: storage.Message{
			Data:       data,
			Attributes: attributes,
		},
		MessageType: "cancel",
	}

	t.Run("publish message", func(t *testing.T) {
		ticker, mocks := tickerWithMocks()

		mocks.timeSource.On("Now").Return(now)
		mocks.router.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := ticker.HandleMessage(ctx, msg)
		assert.NoError(t, err)

		mocks.router.AssertCalled(t, "Publish", mock.Anything, "cancel", data, attributes)
	})
}

type tickerMocks struct {
	router     *schedulermocks.Router
	storage    *storagemocks.SchedulerStorage
	timeSource *schedulermocks.TimeSource
}

func tickerWithMocks() (*scheduler.Ticker, tickerMocks) {
	m := tickerMocks{
		router:     &schedulermocks.Router{},
		storage:    &storagemocks.SchedulerStorage{},
		timeSource: &schedulermocks.TimeSource{},
	}

	ticker := scheduler.NewTicker(scheduler.TickerConfig{
		Router:     m.router,
		Storage:    m.storage,
		TimeSource: m.timeSource,
	})

	return ticker, m
}

func assertLogs(t *testing.T, expectedLogs []log, actualLogs []*logrus.Entry) {
	require.Equal(t, len(expectedLogs), len(actualLogs))
	for i, log := range actualLogs {
		assert.Equal(t, expectedLogs[i].msg, log.Message)
		assert.Equal(t, expectedLogs[i].level, log.Level)
	}
}

type log struct {
	msg   string
	level logrus.Level
}
