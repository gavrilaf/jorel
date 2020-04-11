// Code generated by mockery v1.0.0. DO NOT EDIT.

package msgqueuemocks

import (
	context "context"

	msgqueue "github.com/gavrilaf/jorel/pkg/msgqueue"
	mock "github.com/stretchr/testify/mock"
)

// PublisherFactory is an autogenerated mock type for the PublisherFactory type
type PublisherFactory struct {
	mock.Mock
}

// NewPublisher provides a mock function with given fields: ctx, projectID, topicID
func (_m *PublisherFactory) NewPublisher(ctx context.Context, projectID string, topicID string) (msgqueue.Publisher, error) {
	ret := _m.Called(ctx, projectID, topicID)

	var r0 msgqueue.Publisher
	if rf, ok := ret.Get(0).(func(context.Context, string, string) msgqueue.Publisher); ok {
		r0 = rf(ctx, projectID, topicID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(msgqueue.Publisher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, projectID, topicID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}