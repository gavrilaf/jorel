package testdata

import (
	"fmt"
	"time"
)

const (
	AcceptableDelta = 2 * time.Second
)

type Message struct {
	ID               string
	Created          time.Time
	ScheduleDuration time.Duration
}

func (m *Message) String() string {
	return fmt.Sprintf("{ID=%s, Create=%v, Duration=%v}", m.ID, m.Created, m.ScheduleDuration)
}

func (m* Message) Check() (bool, time.Duration) {
	diff := abs(time.Now().Sub(m.Created.Add(m.ScheduleDuration)))
	return diff <= AcceptableDelta, diff
}

func abs(a time.Duration) time.Duration {
	if a >= 0 {
		return a
	}
	return -a
}