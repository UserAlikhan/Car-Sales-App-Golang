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

func UpdateCarModel(carModel *models.CarModelsModel) (*models.CarModelsModel, error) {
	result := database.DB.Model(&models.CarModelsModel{}).Where("id = ?", carModel.ID).Updates(carModel)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedCarModel *models.CarModelsModel
	result = database.DB.First(&updatedCarModel, carModel.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return updatedCarModel, nil
}

func DeleteCarModel(ID uint) error {
	return database.DB.Unscoped().Delete(&models.CarModelsModel{}, ID).Error
}
