package msgqueue

import (
	"fmt"
	"strconv"
)

const (
	delayInSecondsKey = "jor-el-delay"
)

type MsgAttributes struct {
	DelayInSeconds int
	Original       map[string]string
}

func (m MsgAttributes) GetAttributes() map[string]string {
	attributes := make(map[string]string)

	for k, v := range m.Original {
		attributes[k] = v
	}

	attributes[delayInSecondsKey] = strconv.Itoa(m.DelayInSeconds)

	return attributes
}

func NewMsgAttributes(attributes map[string]string) (MsgAttributes, error) {
	var delay int
	var err error

	if strDelay, ok := attributes[delayInSecondsKey]; ok {
		delay, err = strconv.Atoi(strDelay)
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
		DelayInSeconds: delay,
		Original:       other,
	}, nil
}
