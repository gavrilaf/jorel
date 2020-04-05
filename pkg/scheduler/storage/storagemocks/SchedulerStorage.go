// Code generated by mockery v1.0.0. DO NOT EDIT.

package storagemocks

import (
	context "context"

	storage "github.com/gavrilaf/dyson/pkg/scheduler/storage"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// SchedulerStorage is an autogenerated mock type for the SchedulerStorage type
type SchedulerStorage struct {
	mock.Mock
}

// GetMessages provides a mock function with given fields: ctx, olderThan, limit
func (_m *SchedulerStorage) GetMessages(ctx context.Context, olderThan time.Time, limit int) ([]storage.ScheduledMessage, error) {
	ret := _m.Called(ctx, olderThan, limit)

	var r0 []storage.ScheduledMessage
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, int) []storage.ScheduledMessage); ok {
		r0 = rf(ctx, olderThan, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]storage.ScheduledMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, int) error); ok {
		r1 = rf(ctx, olderThan, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, scheduledTime, msg
func (_m *SchedulerStorage) Save(ctx context.Context, scheduledTime time.Time, msg storage.Message) error {
	ret := _m.Called(ctx, scheduledTime, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, storage.Message) error); ok {
		r0 = rf(ctx, scheduledTime, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
