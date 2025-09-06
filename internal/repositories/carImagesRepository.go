package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarImage(path string, carPostID uint) error {
	carImage := models.CarImagesModel{
		Path:      path,
		CarPostID: carPostID,
	}

	return database.DB.Create(&carImage).Error
}

func DeleteCarImageDBRecord(ID uint) error {
	return database.DB.Unscoped().Delete(&models.CarImagesModel{}, ID).Error
}
