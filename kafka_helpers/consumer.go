package kafka_helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"internship_project/elasticsearch_helpers"
	"log"
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

		var jsonMessage map[string]interface{}
		err = json.Unmarshal(m.Value, &jsonMessage)
		if err != nil {
			fmt.Printf("Unable to parse Kafka message")
			continue
		}

		if jsonMessage["operation"] == OperationEnumString(Created) || jsonMessage["operation"] == OperationEnumString(Updated) {
			product, err := json.Marshal(jsonMessage["product"])
			if err != nil {
				fmt.Println(err)
				continue
			}

			consumer.EsClient.IndexDocument(string(m.Key), string(product))
		} else if jsonMessage["operation"] == OperationEnumString(Deleted) {
			consumer.EsClient.DeleteDocument(string(m.Key))
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
