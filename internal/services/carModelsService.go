package services

import (
	"car_sales/internal/models"
	"car_sales/internal/repositories"
)

func CreateCarModel(carModel *models.CarModelsModel) error {
	return repositories.CreateCarModel(carModel)
}

func GetAllCarModels() ([]*models.CarModelsModel, error) {
	return repositories.GetAllCarModels()
}

func GetCarModelsByCarBrandID(carBrandID int) ([]*models.CarModelsModel, error) {
	return repositories.GetCarModelsByCarBrandID(uint(carBrandID))
}
