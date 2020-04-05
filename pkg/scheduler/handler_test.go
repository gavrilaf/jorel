package scheduler_test

import (
	"context"
	"github.com/gavrilaf/dyson/pkg/scheduler/schedulermocks"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/dyson/pkg/msgqueue/msgqueuemocks"
	"github.com/gavrilaf/dyson/pkg/scheduler"
	"github.com/gavrilaf/dyson/pkg/scheduler/storage/storagemocks"
)

func TestHandlerReceive(t *testing.T) {
	ctx := context.Background()
	data := []byte("data")

	publishResult := &msgqueuemocks.PublishResult{}
	publishResult.On("GetMessageID", mock.Anything).Return("1", nil)

	t.Run("no delay attribute", func(t *testing.T) {
		subject, _ := subjectWithMocks()

		err := subject.Receive(ctx, data, map[string]string{})
		assert.Error(t, err)
	})

	t.Run("immediately resend empty attributes", func(t *testing.T) {
		subject, m := subjectWithMocks()

		m.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := subject.Receive(ctx, data, map[string]string{"jor-el-delay": "0"})
		assert.NoError(t, err)

		var empty map[string]string
		m.publisher.AssertCalled(t, "Publish", mock.Anything, data, empty)
	})

	t.Run("immediately resend with attributes", func(t *testing.T) {
		subject, m := subjectWithMocks()

		m.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := subject.Receive(ctx, data, map[string]string{"jor-el-delay": "0", "one": "two"})
		assert.NoError(t, err)

		m.publisher.AssertCalled(t, "Publish", mock.Anything, data, map[string]string{"one": "two"})
	})
}

type mocks struct {
	publisher *msgqueuemocks.Publisher
	storage   *storagemocks.SchedulerStorage
	timeSource *schedulermocks.TimeSource
}

func subjectWithMocks() (*scheduler.Handler, mocks) {
	m := mocks{
		publisher: &msgqueuemocks.Publisher{},
		storage:   &storagemocks.SchedulerStorage{},
		timeSource: &schedulermocks.TimeSource{},
	}

	subject := scheduler.NewHandler(scheduler.HandlerConfig{
		Publisher: m.publisher,
	})

	return subject, m
}
