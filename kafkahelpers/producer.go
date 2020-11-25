package kafkahelpers

import (
	"context"
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func (producer *KafkaProducer) WriteMessage(topicName string, message string) error {
	kafkaMessage := kafka.Message{
		Key:   []byte(uuid.NewV4().String()),
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
