package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// Define the client configuration with your credentials and Elasticsearch address
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://172.19.255.204:9200", // Replace with your Elasticsearch address
		},
		Username: "aes",     // Replace with your username
		Password: "aes1234", // Replace with your password
	}

	// Initialize the Elasticsearch client
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	now := time.Now().UTC()
	startTime := now.Add(-48 * time.Hour).Format(time.RFC3339)
	endTime := now.Format(time.RFC3339)

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("kube_containers"),
		es.Search.WithBody(strings.NewReader(string(LogQuery("unknown", "aas-ns", "aaa-sql-1-db", startTime, endTime)))), // Convert []byte to string
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// Decode the response
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// Retrieve the total number of matched logs
	if hits, ok := response["hits"].(map[string]interface{}); ok {
		total := hits["total"].(map[string]interface{})["value"].(float64)
		fmt.Printf("Total Doc: %d\n", int(total))
		fmt.Println(len(hits["hits"].([]interface{})))
		// Iterate over the logs and print them
		for _, hit := range hits["hits"].([]interface{}) {
			logEntry := hit.(map[string]interface{})
			source := logEntry["_source"].(map[string]interface{})
			timestamp := source["@timestamp"]
			message := source["message"]
			fmt.Printf("Timestamp: %s, Message: %s\n", timestamp, message)
		}
	} else {
		fmt.Println("No logs found")
	}

	// // Perform the count request
	// res, err := es.Count(
	// 	es.Count.WithContext(context.Background()),
	// 	es.Count.WithIndex("kube_containers"),
	// 	es.Count.WithBody(strings.NewReader(string(LogQuery("warning", "aas-ns", "aaa-sql-1-db")))),
	// )
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// }
	// defer res.Body.Close()

	// // Decode the response
	// var response map[string]interface{}
	// if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
	// 	log.Fatalf("Error parsing the response body: %s", err)
	// }

	// // Print the total count of matched documents
	// if count, ok := response["count"].(float64); ok {
	// 	fmt.Printf("count: %d\n", int(count))
	// } else {
	// 	fmt.Println("Error retrieving count")
	// }

}
