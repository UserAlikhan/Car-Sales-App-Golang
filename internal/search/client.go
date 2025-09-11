package search

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
)

var ES *elasticsearch.Client

func InitElasticSearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Errorf("Error creating ES client %s ", err)
	}

	// simple ping / health check
	_, err = client.Info()
	if err != nil {
		fmt.Errorf("elasticsearch is not reachable %v", err)
	}

	ES = client
	fmt.Println("Elasticsearch was initialized.")
}
