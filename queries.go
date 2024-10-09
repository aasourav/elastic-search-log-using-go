package main

import (
	"encoding/json"
	"log"
)

func buildMatchPhrases(phrases []string) []map[string]interface{} {
	var clauses []map[string]interface{}

	for _, phrase := range phrases {
		clauses = append(clauses, map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"message": phrase,
			},
		})
	}

	return clauses
}

func getMatchPhrases(logType string) []string {
	phrases := map[string][]string{
		"error":   {"error", "panic", "critical", "emergency", "alert", "security"},
		"info":    {"notice", "info", "system", "transaction", "system", "profile", "access"},
		"warning": {"warning", "notice", "diagnostic"},
		"debug":   {"debug", "trace", "verbose"},
		"unknown": {"error", "panic", "critical", "emergency", "alert", "security", "notice", "info", "system", "transaction", "system", "profile", "access", "warning", "notice", "diagnostic", "debug", "trace", "verbose"},
	}

	return phrases[logType]
}

func LogQuery(logType string, namespace string, deployment string, startTime string, endTime string) []byte {
	matchPhrases := getMatchPhrases(logType)

	shouldQuery := map[string]interface{}{
		"from": 0,
		"size": 10000,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"kubernetes.namespace_labels.kubernetes_io/metadata_name.keyword": namespace,
						},
					},
					{
						"term": map[string]interface{}{
							"kubernetes.labels.app_kubernetes_io/instance.keyword": deployment,
						},
					},
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte":    startTime,                   // Greater than or equal to startTime
								"lte":    endTime,                     // Less than or equal to endTime
								"format": "strict_date_optional_time", // Timestamp format
							},
						},
					},
				},
				"should":               buildMatchPhrases(matchPhrases),
				"minimum_should_match": 1,
			},
		},
	}

	allQuery := map[string]interface{}{
		"from": 0,
		"size": 10000,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"kubernetes.namespace_labels.kubernetes_io/metadata_name.keyword": namespace,
						},
					},
					{
						"term": map[string]interface{}{
							"kubernetes.labels.app_kubernetes_io/instance.keyword": deployment,
						},
					},
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte":    startTime,                   // Greater than or equal to startTime
								"lte":    endTime,                     // Less than or equal to endTime
								"format": "strict_date_optional_time", // Timestamp format
							},
						},
					},
				},
			},
		},
	}

	mustNoQuery := map[string]interface{}{
		"from": 0,
		"size": 10000,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"kubernetes.namespace_labels.kubernetes_io/metadata_name.keyword": namespace,
						},
					},
					{
						"term": map[string]interface{}{
							"kubernetes.labels.app_kubernetes_io/instance.keyword": deployment,
						},
					},
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte":    startTime,                   // Greater than or equal to startTime
								"lte":    endTime,                     // Less than or equal to endTime
								"format": "strict_date_optional_time", // Timestamp format
							},
						},
					},
				},
				"must_not": buildMatchPhrases(matchPhrases),
			},
		},
	}

	query := shouldQuery
	if logType == "unknown" {
		query = mustNoQuery
	} else if logType == "" {
		query = allQuery
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error marshalling query: %s", err)
	}

	return queryJSON
}
