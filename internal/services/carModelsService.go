package services

import (
	"car_sales/internal/cache"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"encoding/json"
	"fmt"
	"time"
)

func CreateCarModel(carModel *models.CarModelsModel) error {
	return repositories.CreateCarModel(carModel)
}

func GetAllCarModels() ([]*models.CarModelsModel, error) {
	return repositories.GetAllCarModels()
}

func GetCarModelsByCarBrandID(carBrandID int) ([]*models.CarModelsModel, error) {
	carModelByCarIDCacheKey := fmt.Sprintf("carbrand:%d:carmodels", carBrandID)
	var carModels []*models.CarModelsModel

	// 1. Try to get car models based on carbrand id from the Redis cache first
	cachedData, err := cache.GetCache(carModelByCarIDCacheKey)
	if err == nil && cachedData != "" {
		// Unmarshal the json
		err = json.Unmarshal([]byte(cachedData), &carModels)
		if err != nil {
			// if we cannot unmarshal the json (we got an error)
			// there is something wrong with this json, so delete the cache
			cache.DeleteCache(carModelByCarIDCacheKey)
		} else {
			// If there is no error return data from the cache
			return carModels, nil
		}
	}

	// 2. If there is no appropriate data in the Redis cache, fetch data from the DB
	carModel, err := repositories.GetCarModelsByCarBrandID(uint(carBrandID))
	if err != nil {
		return nil, err
	}

	// save data to the Redis cache
	marshaledData, err := json.Marshal(carModel)
	if err == nil {
		cache.SetCache(carModelByCarIDCacheKey, string(marshaledData), 24*time.Hour)
	}

	return carModel, nil
}

func UpdateCarModel(carModel *models.CarModelsModel) (*models.CarModelsModel, error) {
	return repositories.UpdateCarModel(carModel)
}

func DeleteCarModel(ID uint) error {
	return repositories.DeleteCarModel(ID)
}
