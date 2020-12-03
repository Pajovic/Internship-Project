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

		m, err := consumer.Reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Error while fetching message")
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		s := strings.Split(string(m.Value), " ")

		if len(s) == 3 {
			if s[0] == "CREATED" || s[0] == "UPDATED" {
				go consumer.EsClient.IndexDocument(string(m.Key), s[2])
			} else if s[0] == "DELETED" {
				go consumer.EsClient.DeleteDocument(string(m.Key))
			}
		} else {
			log.Printf("Failed to parse Kafka message")
		}

		err = consumer.Reader.CommitMessages(context.Background(), m)
		if err != nil {
			log.Printf("Failed to commit message")
		}

	}

	if err := consumer.Reader.Close(); err != nil {
		log.Printf("failed to close reader: ", err)
	}
}
