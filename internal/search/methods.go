package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

// Creates or replaces document in Elastic Search Index
func CreateCarPostES(ctx context.Context, doc CarPostDoc) error {
	if ES == nil {
		return fmt.Errorf("es is not initializedS")
	}

	body, _ := json.Marshal(doc)

	req := esapi.IndexRequest{
		Index:      CarPostIndex,
		DocumentID: fmt.Sprintf("%d", doc.ID),
		Body:       bytes.NewReader(body),
		Refresh:    "false",
	}

	res, err := req.Do(ctx, ES)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index error: %s", res.String())
	}

	return nil

}

// This method updates carpost data in the Elastic Search
func UpdateCarPost(ctx context.Context, id uint, fields map[string]interface{}) error {
	if ES == nil {
		return fmt.Errorf("es is not initialized.")
	}

	payload := map[string]interface{}{"doc": fields}

	body, _ := json.Marshal(payload)

	req := esapi.UpdateRequest{
		Index:      CarPostIndex,
		DocumentID: fmt.Sprintf("%d", id),
		Body:       bytes.NewReader(body),
		Refresh:    "false",
	}

	res, err := req.Do(ctx, ES)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update error: %s", res.String())
	}

	return nil
}

// This function deletes car post in the Elastic Search
func DeleteCarPost(ctx context.Context, id uint) error {
	if ES == nil {
		return fmt.Errorf("es in not initialized.")
	}

	req := esapi.DeleteRequest{
		Index:      CarPostIndex,
		DocumentID: fmt.Sprintf("%d", id),
		Refresh:    "false",
	}

	res, err := req.Do(ctx, ES)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete error: %s", res.String())
	}

	return nil
}
