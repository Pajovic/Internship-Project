package kafka_helpers

import (
	"github.com/segmentio/kafka-go"
	"internship_project/elasticsearch_helpers"
)

func NewConsumer(topicName string, EsClient elasticsearch_helpers.ElasticsearchClient) KafkaConsumer {
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

func GetWriter(topicName string) *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	})

	return w
}
