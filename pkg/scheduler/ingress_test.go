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

func TestHandlerReceive(t *testing.T) {
	ctx := context.Background()
	data := []byte("data")

	now := time.Now()

	publishResult := &msgqueuemocks.PublishResult{}
	publishResult.On("GetMessageID", mock.Anything).Return("1", nil)

	t.Run("no delay attribute", func(t *testing.T) {
		ingress, _ := ingressWithMocks()

		err := ingress.Receive(ctx, data, map[string]string{})
		assert.Error(t, err)
	})

	t.Run("immediately resend empty attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := ingress.Receive(ctx, data, map[string]string{"jor-el-delay": "0"})
		assert.NoError(t, err)

		var empty map[string]string
		m.publisher.AssertCalled(t, "Publish", mock.Anything, data, empty)
	})

	t.Run("immediately resend with attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := ingress.Receive(ctx, data, map[string]string{"jor-el-delay": "0", "one": "two"})
		assert.NoError(t, err)

		m.publisher.AssertCalled(t, "Publish", mock.Anything, data, map[string]string{"one": "two"})
	})

	t.Run("when publisher fails", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf(""))

		err := ingress.Receive(ctx, data, map[string]string{"jor-el-delay": "0", "one": "two"})
		assert.Error(t, err)
	})

	t.Run("save with empty attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.timeSource.On("Now").Return(now)
		m.storage.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := ingress.Receive(ctx, data, map[string]string{"jor-el-delay": "5"})
		assert.NoError(t, err)

		expectedTime := now.Add(5 * time.Second)
		expectedMessage := storage.Message{
			Data: data,
		}

		m.storage.AssertCalled(t, "Save", mock.Anything, expectedTime, expectedMessage)
	})

	t.Run("save with attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.timeSource.On("Now").Return(now)
		m.storage.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := ingress.Receive(ctx, data, map[string]string{"jor-el-delay": "20", "one": "two"})
		assert.NoError(t, err)

		expectedTime := now.Add(20 * time.Second)
		expectedMessage := storage.Message{
			Data:       data,
			Attributes: map[string]string{"one": "two"},
		}

		m.storage.AssertCalled(t, "Save", mock.Anything, expectedTime, expectedMessage)
	})

	t.Run("when storage fails", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.timeSource.On("Now").Return(now)
		m.storage.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf(""))

		err := ingress.Receive(ctx, data, map[string]string{"jor-el-delay": "5"})
		assert.Error(t, err)
	})
}

type ingressMocks struct {
	publisher  *msgqueuemocks.Publisher
	storage    *storagemocks.SchedulerStorage
	timeSource *schedulermocks.TimeSource
}

func ingressWithMocks() (*scheduler.Ingress, ingressMocks) {
	m := ingressMocks{
		publisher:  &msgqueuemocks.Publisher{},
		storage:    &storagemocks.SchedulerStorage{},
		timeSource: &schedulermocks.TimeSource{},
	}

	ingress := scheduler.NewIngress(scheduler.IngressConfig{
		Publisher:  m.publisher,
		Storage:    m.storage,
		TimeSource: m.timeSource,
	})

	return ingress, m
}
