package elasticsearch_helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
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

func (esclient *ElasticsearchClient) IndexDocument(id string, body string) {
	req := esapi.IndexRequest{
		Index:      "product",
		DocumentID: id,
		Body:       strings.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), esclient.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), id)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}

func (esclient *ElasticsearchClient) DeleteDocument(id string) {
	req := esapi.DeleteRequest{
		Index: "product",
		DocumentID: id,
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), esclient.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error deleting document ID=%s", res.Status(), id)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}