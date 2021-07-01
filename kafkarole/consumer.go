package kafkarole

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewConsumer(serverAddr, groupId string, topics []string) (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": serverAddr,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	return c, nil
}
