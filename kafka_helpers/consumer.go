package kafka_helpers

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"internship_project/elasticsearch_helpers"
	"log"
	"strings"
)

type KafkaConsumer struct {
	Reader   *kafka.Reader
	EsClient elasticsearch_helpers.ElasticsearchClient
}

func (consumer *KafkaConsumer) Consume() {
	fmt.Println("KafkaConsumer is ready to consume new messages.")
	for {

		m, err := consumer.Reader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		s := strings.Split(string(m.Value), " ")

		if s[0] == "CREATED" || s[0] == "UPDATED" {
			go consumer.EsClient.IndexDocument(string(m.Key), s[2])
		} else if s[0] == "DELETED" {
			go consumer.EsClient.DeleteDocument(string(m.Key))
		}
	}

	if err := consumer.Reader.Close(); err != nil {
		log.Fatal("failed to close reader: ", err)
	}
}
