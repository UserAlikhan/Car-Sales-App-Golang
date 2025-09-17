package utils

import (
	"car_sales/internal/models"
	"car_sales/internal/search"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// this function gets all the car posts and inserts them into elastic search one by one
func InsertDataIntoElasticSearch() error {
	fmt.Println("\nINSERTING DATA INTO ELASTIC SEARCH\n")
	baseUrl := "http://localhost:3000/"
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OCwidXNlcm5hbWUiOiIiLCJlbWFpbCI6ImNoZWxAbWFpbC5ydSIsImlzQWRtaW4iOnRydWUsImV4cCI6MTc2MDcxMDc0NH0.8vFFp4cwHdfFQ3r_bprBuk6G0_KVHhobvGtQUw3SSIo"

	// call carPosts endpoin to get all car posts
	req, err := http.NewRequest("GET", baseUrl+"carPost/getAllCarPosts", nil)
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}

	// add authorization header
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error getting data from the database: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-200 status: %d", response.StatusCode)
	}

	// read and unmarshal response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var carPosts []models.CarPostsModel
	if err := json.Unmarshal(body, &carPosts); err != nil {
		return fmt.Errorf("Failed to unmarshal response: %w", err)
	}

	// insert each car post into Elastic search
	for i := 0; i < len(carPosts); i++ {

		var carPostDoc = search.CarPostDoc{
			ID:            carPosts[i].ID,
			Year:          carPosts[i].Year,
			Description:   carPosts[i].Description,
			Mileage:       carPosts[i].Mileage,
			Price:         carPosts[i].Price,
			ExteriorColor: carPosts[i].ExteriorColor,
			InteriorColor: carPosts[i].InteriorColor,
			Brand:         carPosts[i].CarModel.CarBrand.Name,
			Model:         carPosts[i].CarModel.Name,
			Address:       carPosts[i].Address,
		}

		err = search.CreateCarPostES(context.Background(), carPostDoc)
		if err != nil {
			return fmt.Errorf("Failed to insert data into Elastic Search %s", err.Error())
		}
	}

	return nil
}
