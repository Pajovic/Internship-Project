package elasticsearch_helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

func GetElasticsearchClient() ElasticsearchClient {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error getting client: %s", err)
	}
	return ElasticsearchClient{
		client: es,
	}
}

func (esclient *ElasticsearchClient) SearchDocument(term string) []byte {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"name": "*" + term + "*",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := esclient.client.Search(
		esclient.client.Search.WithContext(context.Background()),
		esclient.client.Search.WithIndex("product"),
		esclient.client.Search.WithBody(&buf),
		esclient.client.Search.WithTrackTotalHits(true),
		esclient.client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	var final []map[string]interface{}

	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		final = append(final, hit.(map[string]interface{})["_source"].(map[string]interface{}))
	}

	json, _ := json.MarshalIndent(final, "", "    ")

	return json
}