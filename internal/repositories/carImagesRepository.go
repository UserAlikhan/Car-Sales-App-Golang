package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarImage(path string, carPostId uint) error {
	carImage := models.CarImagesModel{
		Path:      path,
		CarPostId: carPostId,
	}

	return database.DB.Create(&carImage).Error
}
