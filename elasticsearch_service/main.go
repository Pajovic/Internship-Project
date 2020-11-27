package main

import (
	elkafka "internship_project/elasticsearch_service/kafkaHelpers"
	"internship_project/elasticsearch_service/elasticsearchHelpers"
	"sync"
)

var (
	r  map[string]interface{}
	wg sync.WaitGroup
)

func main() {
	EsClient := elasticsearchHelpers.GetElasticsearchClient()

	kafkaConsumer := elkafka.NewConsumer("ava-internship", EsClient)
	kafkaConsumer.Consume()
}
