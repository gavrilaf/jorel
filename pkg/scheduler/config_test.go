package scheduler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/jorel/pkg/scheduler"
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
