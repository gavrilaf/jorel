package scheduler

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type AggregationType int

const (
	AggregationTypeNo     = iota
	AggregationTypeSkip   = iota
	AggregationTypeAppend = iota
)

type EgressTopicConfig struct {
	Name        string
	Aggregation AggregationType
}

type Config struct {
	ProjectID           string
	IngressSubscription string
	DefaultEgress       EgressTopicConfig
	Routing             map[string]EgressTopicConfig
}

func ReadConfig(buffer []byte) (Config, error) {

	type YamlRouteConfig struct {
		MessageType string `yaml:"message-type"`
		TopicName   string `yaml:"topic-name"`
		Aggregation string `yaml:"aggregation"`
	}

	type YamlConfig struct {
		ProjectID           string            `yaml:"project-id"`
		IngressSubscription string            `yaml:"ingress-subscription"`
		DefaultEgress       YamlRouteConfig   `yaml:"default"`
		Routing             []YamlRouteConfig `yaml:"routing"`
	}

	var yamlConfig YamlConfig
	err := yaml.Unmarshal(buffer, &yamlConfig)
	if err != nil {
		return Config{}, err
	}

	if yamlConfig.ProjectID == "" {
		return Config{}, fmt.Errorf("empty project id")
	}

	if yamlConfig.IngressSubscription == "" {
		return Config{}, fmt.Errorf("empty ingress subscription")
	}

	if yamlConfig.DefaultEgress.TopicName == "" {
		return Config{}, fmt.Errorf("empty default egress topic name")
	}

	routing := make(map[string]EgressTopicConfig)
	for _, r := range yamlConfig.Routing {
		routing[r.MessageType] = EgressTopicConfig{
			Name: r.TopicName,
		}
	}

	return Config{
		ProjectID:           yamlConfig.ProjectID,
		IngressSubscription: yamlConfig.IngressSubscription,
		DefaultEgress: EgressTopicConfig{
			Name: yamlConfig.DefaultEgress.TopicName,
		},
		Routing: routing,
	}, nil
}
