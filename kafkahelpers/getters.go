package kafkahelpers

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func GetConnection() *kafka.Conn {
	topic := "ava-internship"
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return conn
}

func NewConsumer(topicName string) KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topicName,
		Partition: 0,
		MinBytes:  10e2, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	r.SetOffset(0)

	consumer := KafkaConsumer{
		Reader: r,
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
