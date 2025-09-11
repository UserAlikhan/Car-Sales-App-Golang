package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarModel(carModel *models.CarModelsModel) error {
	return database.DB.Create(carModel).Error
}

func GetAllCarModels() ([]*models.CarModelsModel, error) {
	var carModels []*models.CarModelsModel

	result := database.DB.Find(&carModels)
	if result.Error != nil {
		return nil, result.Error
	}

	return carModels, nil
}

func GetCarModelsByCarBrandID(carBrandID uint) ([]*models.CarModelsModel, error) {
	var carModels []*models.CarModelsModel

	result := database.DB.Where(&models.CarModelsModel{CarBrandID: carBrandID}).Find(&carModels)
	if result.Error != nil {
		return nil, result.Error
	}
	return carModels, nil
}
