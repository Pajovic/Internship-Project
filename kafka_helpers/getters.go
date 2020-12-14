package kafka_helpers

import (
	"internship_project/elasticsearch_helpers"
	"time"

	"github.com/segmentio/kafka-go"
)

func NewConsumer(topicName string, address string, groupId string, EsClient elasticsearch_helpers.ElasticsearchClient, miliseconds int) KafkaConsumer {
	r := GetReader(topicName, address, groupId, miliseconds)

	r.SetOffset(kafka.LastOffset)

	consumer := KafkaConsumer{
		Reader:   r,
		EsClient: EsClient,
	}

	return consumer
}

func GetReader(topicName string, address string, groupId string, miliseconds int) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		GroupID:   groupId,
		Topic:     topicName,
		Partition: 0,
		MinBytes:  10e2, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   time.Duration(miliseconds) * time.Millisecond,
	})

	return r
}

func GetWriter(topicName string) *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	})

	return w
}

func GetRetryHandler(retryTopic string, mainTopic string, address string, groupId string, miliseconds int) *RetryHandler {
	handler := &RetryHandler{
		Reader: GetReader(retryTopic, address, groupId, miliseconds),
		Writer: &KafkaProducer{
			Writer: GetWriter(mainTopic),
		},
	}

	return handler
}
