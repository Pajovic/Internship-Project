package kafka_helpers

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func (producer *KafkaProducer) WriteMessage(message string, id string) error {
	kafkaMessage := kafka.Message{
		Key:   []byte(id),
		Value: []byte(message),
	}

	err := producer.Writer.WriteMessages(context.Background(), kafkaMessage)
	fmt.Println(fmt.Sprintf("Written message to %s topic", producer.Writer.Topic))

	if err != nil {
		log.Println("failed to write messages:", err)
	}

	return err
}
