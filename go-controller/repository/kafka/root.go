package kafka

import (
	"controller_server_golang/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	cfg *config.Config

	Consumer *kafka.Consumer
}

func NewKafka(cfg *config.Config) (*Kafka, error) {
	k := &Kafka{cfg: cfg}

	var err error

	if k.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
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

func (k *Kafka) Pool(timeoutMs int) kafka.Event {
	return k.Consumer.Poll(timeoutMs)
}

// kafka consumer가 subscribe 할 토픽 지정
// 토픽 = 특정 키 값에 대해 들어오는 값 수용
func (k *Kafka) RegisterSubTopic(topic string) error {
	if err := k.Consumer.Subscribe(topic, nil); err != nil {
		return err
	} else {
		return nil
	}
}
