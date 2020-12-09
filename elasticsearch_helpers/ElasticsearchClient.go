package elasticsearch_helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"internship_project/models"
	"log"
	"strings"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

func GetElasticsearchClient(address string) ElasticsearchClient {
	cfg := elasticsearch.Config{
		Addresses: []string{
			address,
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

func (esclient *ElasticsearchClient) SearchDocument(term string) ([]byte, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"name": "*" + term + "*",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esclient.client.Search(
		esclient.client.Search.WithContext(context.Background()),
		esclient.client.Search.WithIndex("product"),
		esclient.client.Search.WithBody(&buf),
		esclient.client.Search.WithTrackTotalHits(true),
		esclient.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			return nil, errors.New(fmt.Sprintf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],))
		}
	}

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	var final []models.Product

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		jsonProduct, err := json.Marshal(hit.(map[string]interface{})["_source"])
		if err != nil {
			fmt.Println(err)
			continue
		}

		product := models.Product{}
		if err := json.Unmarshal(jsonProduct, &product); err != nil {
			fmt.Println(err)
			continue
		}

		final = append(final, product)
	}

	json, _ := json.MarshalIndent(final, "", "    ")

	return json, nil
}

func (esclient *ElasticsearchClient) IndexDocument(id string, body string) error {
	req := esapi.IndexRequest{
		Index:      "product",
		DocumentID: id,
		Body:       strings.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), esclient.client)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), id)
		return err
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return err
	} else {
		log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	}
	return nil
}

func (esclient *ElasticsearchClient) DeleteDocument(id string) error {
	req := esapi.DeleteRequest{
		Index: "product",
		DocumentID: id,
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), esclient.client)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		err := fmt.Sprint("[%s] Error deleting document ID=%s", res.Status(), id)
		log.Printf(err)
		return errors.New(err)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			return err
		} else {
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
	return nil
}