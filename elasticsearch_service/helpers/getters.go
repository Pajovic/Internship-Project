package helpers

import (
	"github.com/segmentio/kafka-go"
)

func NewConsumer(topicName string, EsClient ElasticsearchClient) KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topicName,
		Partition: 0,
		MinBytes:  10e2, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	r.SetOffset(kafka.LastOffset)

	consumer := KafkaConsumer{
		Reader: r,
		EsClient: EsClient,
	}

	return consumer
}
