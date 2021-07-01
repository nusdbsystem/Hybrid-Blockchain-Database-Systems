package kafkarole

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewProducer(serverAddr, topic string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": serverAddr})
	if err != nil {
		return nil, err
	}
	return p, nil
}
