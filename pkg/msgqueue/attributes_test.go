package msgqueue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/dyson/pkg/msgqueue"
)

func TestMsgAttributes(t *testing.T) {
	happy := []struct {
		name     string
		msgAttrs msgqueue.MsgAttributes
	}{
		{
			name: "only with delay field",
			msgAttrs: msgqueue.MsgAttributes{
				DelayInSeconds: 10,
			},
		},
		{
			name: "empty delay field",
			msgAttrs: msgqueue.MsgAttributes{},
		},
		{
			name: "with other attributes",
			msgAttrs: msgqueue.MsgAttributes{
				DelayInSeconds: 86400,
				Original:       map[string]string{"one": "two", "three": "four"},
			},
		},
	}

	for _, tt := range happy {
		t.Run(tt.name, func(t *testing.T) {
			attrs := tt.msgAttrs.GetAttributes()
			newMa, err := msgqueue.NewMsgAttributes(attrs)
			assert.NoError(t, err)
			assert.Equal(t, tt.msgAttrs, newMa)
		})
	}

	invalidAttrs := []struct {
		name     string
		attrs map[string]string
	}{
		{
			name: "empty attributes",
			attrs: map[string]string{},
		},
		{
			name: "no delay field",
			attrs: map[string]string{"one": "two"},
		},
		{
			name: "negative delay",
			attrs: map[string]string{"jor-el-delay": "-10"},
		},
		{
			name: "invalid delay",
			attrs: map[string]string{"jor-el-delay": "i'm delay"},
		},
	}

	for _, tt := range invalidAttrs {
		t.Run(tt.name, func(t *testing.T) {
			_, err := msgqueue.NewMsgAttributes(tt.attrs)
			assert.Error(t, err)
		})
	}
}
