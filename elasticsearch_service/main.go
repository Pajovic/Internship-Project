package main

import (
	"internship_project/elasticsearch_service/helpers"
	"sync"
)

var (
	r  map[string]interface{}
	wg sync.WaitGroup
)

func main() {
	EsClient := helpers.GetElasticsearchClient()

	kafkaConsumer := helpers.NewConsumer("ava-internship", EsClient)
	kafkaConsumer.Consume()
}
