package scheduler_test

import (
	"github.com/gavrilaf/dyson/pkg/scheduler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	str := `
---
project-id: test
ingress-subscription: ingress-sub

default:
  topic-name: default-egress
  aggregation: no

routing:
  - message-type: type-one
    topic-name: topic-one
    aggregation: no

  - message-type: type-two
    topic-name: topic-two
    aggregation: skip
`
	cfg, err := scheduler.ParseConfig([]byte(str))
	assert.NoError(t, err)

	expected := scheduler.Config{
		ProjectID:           "test",
		IngressSubscription: "ingress-sub",
		DefaultEgress: scheduler.EgressTopicConfig{
			Name: "default-egress",
		},
		Routing: map[string]scheduler.EgressTopicConfig{
			"type-one": {
				Name: "topic-one",
			},
			"type-two": {
				Name: "topic-two",
			},
		},
	}

	assert.Equal(t, expected, cfg)
}
