package kafkaHelpers

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func (consumer *KafkaConsumer) Consume() {
	go func() {
		for {
			m, err := consumer.Reader.ReadMessage(context.Background())

			if err != nil {
				break
			}

			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		}

		if err := consumer.Reader.Close(); err != nil {
			log.Fatal("failed to close reader: ", err)
		}
	}()
}
