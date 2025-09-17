package main

import (
	"car_sales/internal/search"
	"car_sales/internal/utils"
	"fmt"
	"log"
)

func main() {
	fmt.Println("🚀 Syncing car posts into Elasticsearch...")

	// init elastic seach before calling insert data method
	search.InitElasticSearch()

	if err := utils.InsertDataIntoElasticSearch(); err != nil {
		log.Fatalf("❌ Failed: %v", err)
	}

	fmt.Println("✅ Done!")

}
