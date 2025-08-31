package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarModel(carModel *models.CarModelsModel) error {
	return database.DB.Create(carModel).Error
}
