package kafkahelpers

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func (producer *KafkaProducer) WriteMessage(topicName string, message string, id string) error {
	kafkaMessage := kafka.Message{
		Key:   []byte(id),
		Value: []byte(message),
	}

	err := producer.Writer.WriteMessages(context.Background(), kafkaMessage)

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := producer.Writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	return err
}
