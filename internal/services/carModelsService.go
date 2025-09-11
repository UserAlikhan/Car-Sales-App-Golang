package services

import (
	"car_sales/internal/cache"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"encoding/json"
	"fmt"
	"time"
)

// key to get all car models from the redis cache
var carModelsCacheKey = "allcarmodels"

func CreateCarModel(carModel *models.CarModelsModel) error {
	err := repositories.CreateCarModel(carModel)
	if err != nil {
		return err
	}

	// Delete all car models cache
	cache.DeleteCache(carModelsCacheKey)

	return nil
}

func GetAllCarModels() ([]*models.CarModelsModel, error) {
	var carModels []*models.CarModelsModel

	// 1. Try to get data from the Redis cache first
	cachedData, err := cache.GetCache(carModelsCacheKey)
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &carModels); err == nil {
			return carModels, nil
		} else {
			// if there is data in cache, but we cannot unmarshal it
			// it is probably in a wrong format, so delete the cache
			cache.DeleteCache(carModelsCacheKey)
		}
		return carModels, err
	}

	// 2. If there is no data in cache, fetch data from the DB
	carModels, err = repositories.GetAllCarModels()
	if err != nil {
		return nil, err
	}

	// marshal the data
	marshaledData, err := json.Marshal(carModels)
	if err == nil {
		// save the data from db to cache
		cache.SetCache(carModelsCacheKey, string(marshaledData), 24*time.Hour)
	}

	return carModels, nil
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
	carModels, err = repositories.GetCarModelsByCarBrandID(uint(carBrandID))
	if err != nil {
		return nil, err
	}

	// save data to the Redis cache
	marshaledData, err := json.Marshal(carModels)
	if err == nil {
		cache.SetCache(carModelByCarIDCacheKey, string(marshaledData), 24*time.Hour)
	}

	return carModels, nil
}

func UpdateCarModel(carModel *models.CarModelsModel) (*models.CarModelsModel, error) {
	carModel, err := repositories.UpdateCarModel(carModel)
	if err != nil {
		return nil, err
	}

	// cache name for specific car brand's car models
	carModelByCarPostIDCacheKey := fmt.Sprintf("carbrand:%d:carmodels", carModel.CarBrandID)

	// Delete all car models cache
	cache.DeleteCache(carModelsCacheKey)

	// update specific car model's cache
	marshaledData, err := json.Marshal(carModel)
	if err == nil {
		cache.SetCache(carModelByCarPostIDCacheKey, string(marshaledData), 24*time.Hour)
	}

	return carModel, nil
}

func GetCarModelByID(ID uint) (*models.CarModelsModel, error) {
	return repositories.GetCarModelByID(ID)
}

func DeleteCarModel(ID uint) error {
	// get car model by id
	carModel, err := GetCarModelByID(ID)
	if err != nil {
		return err
	}

	// cache name for specific car brand's car models
	carModelByCarPostIDCacheKey := fmt.Sprintf("carbrand:%d:carmodels", carModel.CarBrandID)

	err = repositories.DeleteCarModel(ID)
	if err != nil {
		return err
	}

	// delete all car models cache
	cache.DeleteCache(carModelsCacheKey)
	// delete a specific car model
	cache.DeleteCache(carModelByCarPostIDCacheKey)

	return nil
}
