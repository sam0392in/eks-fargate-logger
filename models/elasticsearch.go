package models

import (
	"context"
	"eks-fargate-logger/helper"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	logger = helper.Logger()
	client = ElasticsearchSetup()
)

// ElasticsearchSetup : Establish elasticsearch client
// , Reference: https://www.elastic.co/blog/the-go-client-for-elasticsearch-configuration-and-customization
func ElasticsearchSetup() *elasticsearch.Client {
	profile := os.Getenv("ENVIRONMENT") // set environment in k8s values.yaml
	if profile == "" {
		profile = "dev"
	}
	esEndpoint := helper.ReadConfig(fmt.Sprintf("ELASTICSEARCH_%s_ENDPOINT", strings.ToUpper(profile)))
	logger.Info("ES Endpoint: ", esEndpoint)
	var (
		es          *elasticsearch.Client
		clusterURLs = []string{esEndpoint}
		err         error
	)
	cfg := elasticsearch.Config{
		Addresses: clusterURLs,
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error("Error encountered while setting up elasticsearch Client, ERROR: " + err.Error())
	}
	return es
}

// ESInfo : check connection to elasticsearch
func ESInfo() (map[string]interface{}, error) {
	es := ElasticsearchSetup()
	res, err := es.Info()
	if err != nil {
		logger.Error("Error fetching Elasticsearch Info")
		return nil, err
	}
	defer res.Body.Close()

	// Decode the response body into a map
	var info map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		logger.Error("Error decoding es info response")
		return nil, err
	}

	return info, nil
}

func FormatData(indexName string, data map[string]interface{}) esapi.IndexRequest {
	doc, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		logger.Error("Error marshalling data in method Indexation, Error: ", err.Error())
	}

	req := esapi.IndexRequest{
		Index:   indexName,
		Body:    strings.NewReader(string(doc)),
		Refresh: "false",
	}
	return req
}

// There are 2 approaches to send the logs to elasticsearch as given below. Both of them works:
// 1. https://kb.objectrocket.com/elasticsearch/how-to-insert-elasticsearch-documents-into-an-index-using-golang-451
// 2. https://www.elastic.co/blog/the-go-client-for-elasticsearch-working-with-data

// Flush : This function puts the data in elasticsearch endpoint.
func Flush(data []esapi.IndexRequest) {
	for _, req := range data {
		res, err := req.Do(context.Background(), client)
		if err != nil {
			logger.Error("Failed to push data to elasticsearch, ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			logger.Error("ERROR indexing document ERR: %s. status code: %s", res.String(), res.Status())
		} else {

			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				logger.Error("Error deserializing response in indexdata method: ", err.Error())
			}

			result, ok := resMap["result"].(string)
			if !ok || result != "created" {
				logger.Error("document not indexed successfully")
			}

			logger.WithFields(logrus.Fields{
				"doc_id": resMap["_id"],
				"index":  resMap["_index"],
				"shards": resMap["_shards"],
				"result": resMap["result"],
			}).Info("status: ", res.Status())
		}
	}
}
