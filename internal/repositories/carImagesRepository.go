package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarImage(path string, carPostID uint) (*models.CarImagesModel, error) {
	carImage := models.CarImagesModel{
		Path:      path,
		CarPostID: carPostID,
	}

	result := database.DB.Create(&carImage)
	if result.Error != nil {
		return nil, result.Error
	}

	return &carImage, nil
}

func DeleteCarImageDBRecord(ID uint) error {
	return database.DB.Unscoped().Delete(&models.CarImagesModel{}, ID).Error
}

func GetCarImageByID(ID int) (*models.CarImagesModel, error) {
	var carImage *models.CarImagesModel

	result := database.DB.Preload("CarPost").First(&carImage, ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return carImage, nil
}
