package services

import (
	"car_sales/internal/models"
	"car_sales/internal/repositories"
)

func CreateCarModel(carModel *models.CarModelsModel) error {
	return repositories.CreateCarModel(carModel)
}
