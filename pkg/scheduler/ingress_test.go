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

	t.Run("zero delay, empty attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.router.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := ingress.Receive(ctx, data, map[string]string{"delay": "0"})
		assert.NoError(t, err)

		var empty map[string]string
		m.router.AssertCalled(t, "Publish", mock.Anything, "", data, empty)
	})

	t.Run("zero delay, message type, additional attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.router.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

		err := ingress.Receive(ctx, data, map[string]string{"delay": "0", "message-type": "cancel", "one": "two"})
		assert.NoError(t, err)

		m.router.AssertCalled(t, "Publish", mock.Anything, "cancel", data, map[string]string{"one": "two"})
	})

	t.Run("when publisher fails", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.router.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf(""))

		err := ingress.Receive(ctx, data, map[string]string{"delay": "0", "one": "two"})
		assert.Error(t, err)
	})

	t.Run("save with empty attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.timeSource.On("Now").Return(now)
		m.storage.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := ingress.Receive(ctx, data, map[string]string{"delay": "5"})
		assert.NoError(t, err)

		expectedTime := now.Add(5 * time.Second)
		expectedMessage := storage.Message{
			Data: data,
		}

		m.storage.AssertCalled(t, "Save", mock.Anything, expectedTime, "", "", expectedMessage)
	})

	t.Run("save with attributes", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.timeSource.On("Now").Return(now)
		m.storage.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := ingress.Receive(ctx, data, map[string]string{"delay": "20", "message-type": "cancel", "aggregation-id": "123", "one": "two"})
		assert.NoError(t, err)

		expectedTime := now.Add(20 * time.Second)
		expectedMessage := storage.Message{
			Data:       data,
			Attributes: map[string]string{"one": "two"},
		}

		m.storage.AssertCalled(t, "Save", mock.Anything, expectedTime, "cancel", "123", expectedMessage)
	})

	t.Run("when storage fails", func(t *testing.T) {
		ingress, m := ingressWithMocks()

		m.timeSource.On("Now").Return(now)
		m.storage.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf(""))

		err := ingress.Receive(ctx, data, map[string]string{"delay": "5"})
		assert.Error(t, err)
	})
}

type ingressMocks struct {
	router     *schedulermocks.Router
	storage    *storagemocks.SchedulerStorage
	timeSource *schedulermocks.TimeSource
}

func ingressWithMocks() (*scheduler.Ingress, ingressMocks) {
	m := ingressMocks{
		router:     &schedulermocks.Router{},
		storage:    &storagemocks.SchedulerStorage{},
		timeSource: &schedulermocks.TimeSource{},
	}

	ingress := scheduler.NewIngress(scheduler.IngressConfig{
		Router:     m.router,
		Storage:    m.storage,
		TimeSource: m.timeSource,
	})

	return ingress, m
}
