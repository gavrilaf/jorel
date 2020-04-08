package scheduler

import "time"

//go:generate mockery -name TimeSource -outpkg schedulermocks -output ./schedulermocks -dir .
type TimeSource interface {
	Now() time.Time
}

type SystemTime struct{}

func (SystemTime) Now() time.Time {
	return time.Now().UTC()
}
