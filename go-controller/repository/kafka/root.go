package kafka

import (
	"controller_server_golang/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	cfg *config.Config

	consumer *kafka.Consumer
}

func NewKafka(cfg *config.Config) (*Kafka, error) {
	k := &Kafka{cfg: cfg}

	var err error

	if k.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URL,
		"group.id":          cfg.Kafka.GroupId,
		"auto.offset.reset": "latest", // 서버 실행 시 최근 값만 읽음
		"acks":              "all",
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}
