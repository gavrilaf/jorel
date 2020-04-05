package msgqueue

import (
	"fmt"
	"strconv"
	"time"
)

const (
	delayInSecondsKey = "jor-el-delay"
)

type MsgAttributes struct {
	DelayInSeconds time.Duration
	Other          map[string]string
}

func (ma MsgAttributes) GetAttributes() map[string]string {
	attributes := make(map[string]string)

	for k, v := range ma.Other {
		attributes[k] = v
	}

	attributes[delayInSecondsKey] = strconv.FormatInt(int64(ma.DelayInSeconds), 10)

	return attributes
}

func NewMsgAttributes(attributes map[string]string) (MsgAttributes, error) {
	var delay int64
	var err error

	if strDelay, ok := attributes[delayInSecondsKey]; ok {
		delay, err = strconv.ParseInt(strDelay, 10, 64)
		if err != nil {
			return MsgAttributes{}, fmt.Errorf("couldn't parse message delay, %w", err)
		}
		if delay < 0 {
			return MsgAttributes{}, fmt.Errorf("negative delay")
		}
	} else {
		return MsgAttributes{}, fmt.Errorf("couldn't find message delay")
	}

	var other map[string]string
	if len(attributes) > 1 {
		other = make(map[string]string)
		for k, v := range attributes {
			if k != delayInSecondsKey {
				other[k] = v
			}
		}
	}

	return MsgAttributes{
		DelayInSeconds: time.Duration(delay),
		Other:          other,
	}, nil
}
