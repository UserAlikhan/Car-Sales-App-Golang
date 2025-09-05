package services

import (
	"car_sales/internal/models"
	"car_sales/internal/repositories"
)

func CreateCarPost(carPost *models.CarPostsModel) (*models.CarPostsModel, error) {
	return repositories.CreateCarPost(carPost)
}
