package services

import (
	"car_sales/internal/configs"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
)

func CreateCarPost(carPost *models.CarPostsModel) (*models.CarPostsModel, error) {
	return repositories.CreateCarPost(carPost)
}

func GetAllUsersCarPosts(userId uint) ([]*models.CarPostsModel, error) {
	return repositories.GetAllUsersCarPosts(userId)
}

func DeleteCarPost(s3Conf *configs.S3Config, ID uint) error {
	// Check if car post exists first
	carPost, err := repositories.GetCarPostById(ID)
	if err != nil {
		return err
	}

	// call delete car post
	return repositories.DeleteCarPost(s3Conf, carPost, ID)
}
