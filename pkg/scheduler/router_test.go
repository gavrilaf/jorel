package scheduler_test

import (
	"context"
	"fmt"
	"github.com/gavrilaf/jorel/pkg/msgqueue/msgqueuemocks"
	"github.com/gavrilaf/jorel/pkg/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRouterNew(t *testing.T) {
	ctx := context.Background()

	t.Run("happy flow", func(t *testing.T) {
		publisher := &msgqueuemocks.Publisher{}

		factory := &msgqueuemocks.PublisherFactory{}
		factory.On("NewPublisher", mock.Anything, mock.Anything, mock.Anything).Return(publisher, nil)

		config := scheduler.Config{
			ProjectID: "test-prj",
			IngressSubscription: "ingress",
			DefaultEgress: scheduler.EgressTopicConfig{
				Name: "default-egress",
			},
			Routing: map[string]scheduler.EgressTopicConfig{
				"cancel": {
					Name: "cancel-topic",
				},
				"update": {
					Name: "update-topic",
				},
			},
		}

		router, err := scheduler.NewRouter(ctx, config, factory)
		assert.NoError(t, err)
		assert.NotNil(t, router)

		factory.AssertCalled(t, "NewPublisher", mock.Anything, "test-prj", "default-egress")
		factory.AssertCalled(t, "NewPublisher", mock.Anything, "test-prj", "cancel-topic")
		factory.AssertCalled(t, "NewPublisher", mock.Anything, "test-prj", "update-topic")

		factory.AssertNumberOfCalls(t, "NewPublisher", 3)
	})

	t.Run("failed to create default egress topic", func(t *testing.T) {
		factory := &msgqueuemocks.PublisherFactory{}
		factory.On("NewPublisher", mock.Anything, "test-prj", "default-egress").Return(nil, fmt.Errorf(""))

		config := scheduler.Config{
			ProjectID: "test-prj",
			DefaultEgress: scheduler.EgressTopicConfig{
				Name: "default-egress",
			},
		}

		_, err := scheduler.NewRouter(ctx, config, factory)
		assert.Error(t, err)

		factory.AssertExpectations(t)
	})

	t.Run("failed to create roting topic", func(t *testing.T) {
		factory := &msgqueuemocks.PublisherFactory{}
		publisher := &msgqueuemocks.Publisher{}
		factory.On("NewPublisher", mock.Anything, "test-prj", "default-egress").Return(publisher, nil)
		factory.On("NewPublisher", mock.Anything, "test-prj", "cancel-topic").Return(nil, fmt.Errorf(""))

		config := scheduler.Config{
			ProjectID: "test-prj",
			DefaultEgress: scheduler.EgressTopicConfig{
				Name: "default-egress",
			},
			Routing: map[string]scheduler.EgressTopicConfig{
				"cancel": {
					Name: "cancel-topic",
				},
			},
		}

		_, err := scheduler.NewRouter(ctx, config, factory)
		assert.Error(t, err)

		factory.AssertExpectations(t)
	})
}

func TestRouterPublish(t *testing.T) {
	ctx := context.Background()
	router, m := routerWithMocks(t)

	attributes := map[string]string{"one": "two"}

	publishResult := &msgqueuemocks.PublishResult{}

	m.defPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)
	m.updatePublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(publishResult, nil)

	_, err := router.Publish(ctx, "", []byte("empty-message-type"), attributes)
	assert.NoError(t, err)

	_, err = router.Publish(ctx, "update", []byte("update-message"), attributes)
	assert.NoError(t, err)

	m.defPublisher.AssertCalled(t, "Publish", mock.Anything, []byte("empty-message-type"), attributes)
	m.updatePublisher.AssertCalled(t, "Publish", mock.Anything, []byte("update-message"), attributes)
}

type routerMocks struct {
	factory *msgqueuemocks.PublisherFactory
	defPublisher *msgqueuemocks.Publisher
	cancelPublisher *msgqueuemocks.Publisher
	updatePublisher *msgqueuemocks.Publisher
}

func routerWithMocks(t *testing.T) (scheduler.Router, *routerMocks) {
	m := routerMocks{
		factory: &msgqueuemocks.PublisherFactory{},
		defPublisher: &msgqueuemocks.Publisher{},
		cancelPublisher: &msgqueuemocks.Publisher{},
		updatePublisher: &msgqueuemocks.Publisher{},
	}

	m.defPublisher.On("TopicID").Return("default-egress")
	m.cancelPublisher.On("TopicID").Return("cancel-topic")
	m.updatePublisher.On("TopicID").Return("update-topic")

	m.factory.On("NewPublisher", mock.Anything, "test-prj", "default-egress").Return(m.defPublisher, nil)
	m.factory.On("NewPublisher", mock.Anything, "test-prj", "cancel-topic").Return(m.cancelPublisher, nil)
	m.factory.On("NewPublisher", mock.Anything, "test-prj", "update-topic").Return(m.updatePublisher, nil)

	config := scheduler.Config{
		ProjectID: "test-prj",
		IngressSubscription: "ingress",
		DefaultEgress: scheduler.EgressTopicConfig{
			Name: "default-egress",
		},
		Routing: map[string]scheduler.EgressTopicConfig{
			"cancel": {
				Name: "cancel-topic",
			},
			"update": {
				Name: "update-topic",
			},
		},
	}

	router, err := scheduler.NewRouter(context.Background(), config, m.factory)
	assert.NoError(t, err)

	return router, &m
}