package search

import (
	"context"
	"fmt"
	"log"
	"strings"
)

const CarPostIndex = "carposts"

type CarPostDoc struct {
	ID            uint    `json:"id"`
	Year          int     `json:"year"`
	Description   string  `json:"description"`
	Mileage       int     `json:"mileage"`
	Price         float32 `json:"price"`
	ExteriorColor string  `json:"exterior_color"`
	InteriorColor string  `json:"interior_color"`
	Brand         string  `json:"brand,omitempty"`
	Model         string  `json:"model,omitempty"`
	Address       string  `json:"address,omitempty"`
}

// CreateIndex creates the carposts index with a recommended mapping
func CreateIndex(ctx context.Context) error {
	if ES == nil {
		return fmt.Errorf("elasticsearch client was not initialized.")
	}

	// check if index exits first
	res, err := ES.Indices.Exists([]string{CarPostIndex})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// if status code is 200, then index already exists
	if res.StatusCode == 200 {
		return nil
	}

	mapping := `{
		"mappings": {
			"properties": {
				"id":             { "type": "keyword" },
				"year":           { "type": "integer" },
				"description":    { "type": "text", "analyzer": "standard" },
				"mileage":        { "type": "integer" },
				"price":       	  { "type": "float" },
				"exterior_color": { "type": "text" },
				"interior_color": { "type": "text" },
				"brand":          { "type": "keyword" },
				"model":          { "type": "keyword" },
				"address":        { "type": "text" }
			}
		}
	}`

	createRes, err := ES.Indices.Create(CarPostIndex, ES.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return err
	}

	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("failed to create index: %s", createRes.String())
	}

	return nil
}

func DeleteIndex(ctx context.Context) error {
	if ES == nil {
		return fmt.Errorf("elasticsearch client was not initialized.")
	}

	res, err := ES.Indices.Delete([]string{CarPostIndex})
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Index delete failed: %s", res.String())
	}

	log.Println("Index was deleted successfully.")
	return nil
}
