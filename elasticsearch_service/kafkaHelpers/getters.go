package kafkaHelpers

import (
	"github.com/segmentio/kafka-go"
	"internship_project/elasticsearch_service/elasticsearchHelpers"
)

func NewConsumer(topicName string, EsClient elasticsearchHelpers.ElasticsearchClient) KafkaConsumer {
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
