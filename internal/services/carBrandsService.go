package services

import (
	"car_sales/internal/cache"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"encoding/json"
	"fmt"
	"time"
)

/*
Services handle the bussiness logic.
They coordinate multiple repositories (check if user exists + save to DB)
Hide complexities from handler
*/

// GLOBAL VARIABLE
// car brand's cache name
var carBrandCacheKey = fmt.Sprintf("carbrand")

func CreateCarBrand(carBrand *models.CarBrandsModel) error {
	err := repositories.CreateCarBrand(carBrand)
	if err != nil {
		return err
	}

	// if we created a new car brand we need to delete all cache for all car brands
	cache.DeleteCache(carBrandCacheKey)

	return nil
}

func GetAllCarBrands() ([]*models.CarBrandsModel, error) {
	// 1. Try to get data from redis cache first
	var carBrands []*models.CarBrandsModel

	if cached, err := cache.GetCache(carBrandCacheKey); err == nil {
		if err := json.Unmarshal([]byte(cached), &carBrands); err != nil {
			return nil, err
		}

		return carBrands, nil
	}

	// 2. If there is nothing in redis cache, get data from DB
	carBrands, err := repositories.GetAllCarBrands()
	if err != nil {
		return nil, err
	}

	// unmarshal the json
	marsheledData, err := json.Marshal(carBrands)
	if err != nil {
		return nil, err
	}

	// save data to the Redis cache
	cache.SetCache(carBrandCacheKey, string(marsheledData), 24*time.Hour)

	return carBrands, nil
}

func GetCarBrandById(id int) (*models.CarBrandsModel, error) {
	var carBrand *models.CarBrandsModel

	// car brand cache address
	carBrandByIdCacheKey := fmt.Sprintf("carbrand:%d", id)

	// 1. Try to get car brand from the redis cache first
	if cachedData, err := cache.GetCache(carBrandByIdCacheKey); err == nil && cachedData != "" {
		// unmarshal json
		if err := json.Unmarshal([]byte(cachedData), &carBrand); err != nil {
			return nil, err
		}

		// return carBrand from cache
		return carBrand, nil
	}

	// 2. If there is no data in the redis cache, fetch the data from DB
	carBrand, err := repositories.GetCarBrandById(id)
	if err != nil {
		return nil, err
	}

	// marshal the json
	response, _ := json.Marshal(carBrand)

	// save new data to the redis cache
	cache.SetCache(carBrandByIdCacheKey, string(response), 24*time.Hour)

	return carBrand, nil
}

func UpdateCarBrand(carBrand *models.CarBrandsModel) error {
	carBrandByIdCacheKey := fmt.Sprintf("carbrand:%d", carBrand.ID)

	err := repositories.UpdateCarBrand(carBrand)
	if err != nil {
		return err
	}

	// if we updated a car brand we need to delete the cache for all car brand
	cache.DeleteCache(carBrandCacheKey)

	// update cache for a specific car brand
	response, _ := json.Marshal(carBrand)
	cache.SetCache(carBrandByIdCacheKey, string(response), 24*time.Hour)

	return nil
}

func DeleteCarBrand(id int) error {
	err := repositories.DeleteCarBrand(id)
	if err != nil {
		return err
	}

	// if we deleted a car brand we need to delete the cache for all car brands
	cache.DeleteCache(carBrandCacheKey)
	// and for a specific car brand if exists
	carBrandByIdCacheKey := fmt.Sprintf("carbrand:%d", id)
	cache.DeleteCache(carBrandByIdCacheKey)

	return nil
}

func CreateCarBrandWithModels(carBrand *models.CarBrandsModel) error {
	err := repositories.CreateCarBrandWithModels(carBrand)
	if err != nil {
		return err
	}

	// if we created a new car brand we need to delete all cache for all car brands
	cache.DeleteCache(carBrandCacheKey)

	// if we created a new car brand with car models, we need to delete cache for cal models
	if len(carBrand.CarModels) > 0 {
		carModelsCacheID := fmt.Sprintf("carbrand:%d:carmodels", carBrand.ID)
		cache.DeleteCache(carModelsCacheID)
	}

	return nil
}

func SaveCarBrand(carBrand *models.CarBrandsModel) error {
	err := repositories.SaveCarBrand(carBrand)
	if err != nil {
		return err
	}

	// if we created or updated a new car brand we need to delete all cache for all car brands
	cache.DeleteCache(carBrandCacheKey)

	return nil
}
