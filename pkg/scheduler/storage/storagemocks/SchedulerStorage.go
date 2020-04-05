// Code generated by mockery v1.0.0. DO NOT EDIT.

package storagemocks

import (
	time "time"

	storage "github.com/gavrilaf/dyson/pkg/scheduler/storage"
	mock "github.com/stretchr/testify/mock"
)

// SchedulerStorage is an autogenerated mock type for the SchedulerStorage type
type SchedulerStorage struct {
	mock.Mock
}

// GetMessages provides a mock function with given fields: olderThan, limit
func (_m *SchedulerStorage) GetMessages(olderThan time.Time, limit int) ([]storage.ScheduledMessage, error) {
	ret := _m.Called(olderThan, limit)

	var r0 []storage.ScheduledMessage
	if rf, ok := ret.Get(0).(func(time.Time, int) []storage.ScheduledMessage); ok {
		r0 = rf(olderThan, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]storage.ScheduledMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Time, int) error); ok {
		r1 = rf(olderThan, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: scheduledTime, msg
func (_m *SchedulerStorage) Save(scheduledTime time.Time, msg storage.Message) error {
	ret := _m.Called(scheduledTime, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Time, storage.Message) error); ok {
		r0 = rf(scheduledTime, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
