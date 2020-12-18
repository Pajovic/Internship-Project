package kafka_helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"internship_project/elasticsearch_helpers"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader   *kafka.Reader
	EsClient elasticsearch_helpers.ElasticsearchClient
}

func (consumer *KafkaConsumer) Consume() {
	fmt.Println("KafkaConsumer is ready to consume on topic " + consumer.Reader.Stats().Topic)
	retryWriter := GetWriter("retry")
	defer retryWriter.Close()
	retryProducer := KafkaProducer{
		Writer: retryWriter,
	}
	for {
		m, err := consumer.Reader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println("Error while fetching message")
			consumer.resolveError(retryProducer, m)
			continue
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		var jsonMessage map[string]interface{}
		err = json.Unmarshal(m.Value, &jsonMessage)
		if err != nil {
			fmt.Println("Unable to parse Kafka message")
			consumer.resolveError(retryProducer, m)
			continue
		}

		if jsonMessage["operation"] == OperationEnumString(Created) || jsonMessage["operation"] == OperationEnumString(Updated) {
			product, err := json.Marshal(jsonMessage["product"])
			if err != nil {
				fmt.Println(err)
				consumer.resolveError(retryProducer, m)
				continue
			}

			err = consumer.EsClient.IndexDocument(string(m.Key), string(product))
			if err != nil {
				fmt.Println("Error while indexing new Elasticsearch document")
				consumer.resolveError(retryProducer, m)
				continue
			}
		} else if jsonMessage["operation"] == OperationEnumString(Deleted) {
			err = consumer.EsClient.DeleteDocument(string(m.Key))
			if err != nil {
				fmt.Println("Error while deleting Elasticsearch document")
				consumer.resolveError(retryProducer, m)
				continue
			}
		}

		err = consumer.Reader.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println("Failed to commit message")
			continue
		}

	}
}

func (consumer *KafkaConsumer) resolveError(producer KafkaProducer, message kafka.Message) {
	if consumer.Reader.Stats().Topic != "retry" {
		writeToRetry(producer, message)
	}
	err := consumer.Reader.CommitMessages(context.Background(), message)
	if err != nil {
		log.Println("Failed to commit message")
	}
}

func writeToRetry(producer KafkaProducer, message kafka.Message) {
	err := producer.WriteMessage(string(message.Value), string(message.Key))
	if err != nil {
		fmt.Println("Failed to write message to retry topic")
	} else {
		fmt.Println("Successfuly written to retry topic")
	}
}
