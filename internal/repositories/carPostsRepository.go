package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarPost(carPost *models.CarPostsModel) (*models.CarPostsModel, error) {
	return carPost, database.DB.Create(&carPost).Error
}
