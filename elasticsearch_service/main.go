package main

import (
	elkafka "internship_project/elasticsearch_service/elastic_kafkahelpers"
)

func main() {
	kafkaConsumer := elkafka.NewConsumer("ava-internship")
	kafkaConsumer.Consume()
}
