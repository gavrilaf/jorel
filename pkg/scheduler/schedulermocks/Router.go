// Code generated by mockery v1.0.0. DO NOT EDIT.

package schedulermocks

import (
	context "context"

	msgqueue "github.com/gavrilaf/dyson/pkg/msgqueue"
	mock "github.com/stretchr/testify/mock"
)

// Router is an autogenerated mock type for the Router type
type Router struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Router) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: ctx, msgType, data, attributes
func (_m *Router) Publish(ctx context.Context, msgType string, data []byte, attributes map[string]string) (msgqueue.PublishResult, error) {
	ret := _m.Called(ctx, msgType, data, attributes)

	var r0 msgqueue.PublishResult
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, map[string]string) msgqueue.PublishResult); ok {
		r0 = rf(ctx, msgType, data, attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(msgqueue.PublishResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []byte, map[string]string) error); ok {
		r1 = rf(ctx, msgType, data, attributes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}