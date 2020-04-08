package msgqueue

import (
	"fmt"
	"strconv"
)

const (
	delayInSecondsKey = "delay"
	msgTypeKey        = "message-type"
	aggregationIDKey  = "aggregation-id"
)

var predefinedAttributes = map[string]struct{}{
	delayInSecondsKey: {},
	msgTypeKey:        {},
	aggregationIDKey:  {},
}

type MsgAttributes struct {
	DelayInSeconds int
	MessageType    string
	AggregationID  string
	Original       map[string]string
}

func (m MsgAttributes) GetAttributes() map[string]string {
	attributes := make(map[string]string)

	for k, v := range m.Original {
		attributes[k] = v
	}

	attributes[delayInSecondsKey] = strconv.Itoa(m.DelayInSeconds)

	if m.MessageType != "" {
		attributes[msgTypeKey] = m.MessageType
	}

	if m.AggregationID != "" {
		attributes[aggregationIDKey] = m.AggregationID
	}

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

	msgType := attributes[msgTypeKey]
	aggregationID := attributes[aggregationIDKey]

	var other map[string]string
	if len(attributes) > 1 {
		other = make(map[string]string)
		for k, v := range attributes {
			if _, ok := predefinedAttributes[k]; !ok {
				other[k] = v
			}
		}
	}

	return MsgAttributes{
		DelayInSeconds: delay,
		MessageType:    msgType,
		AggregationID:  aggregationID,
		Original:       other,
	}, nil
}
