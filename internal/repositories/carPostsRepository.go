package repositories

import (
	"car_sales/internal/configs"
	"car_sales/internal/database"
	"car_sales/internal/models"
	"car_sales/internal/utils"
	"context"
)

func CreateCarPost(carPost *models.CarPostsModel) (*models.CarPostsModel, error) {
	return carPost, database.DB.Create(&carPost).Error
}

func GetAllUsersCarPosts(userId uint) ([]*models.CarPostsModel, error) {
	var carPosts []*models.CarPostsModel

	result := database.DB.Where(&models.CarPostsModel{SellerID: userId}).Preload("CarModel").Preload("PostImages").Find(&carPosts)
	if result.Error != nil {
		return nil, result.Error
	}

	return carPosts, nil
}

func GetCarPostById(ID uint) (*models.CarPostsModel, error) {
	var carPost models.CarPostsModel

	result := database.DB.First(&carPost, ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &carPost, nil
}

func DeleteCarPost(s3Conf *configs.S3Config, carPost *models.CarPostsModel, ID uint) error {
	// iterate through all images assigned with the post
	for _, img := range carPost.PostImages {
		if err := utils.DeleteFromS3(context.TODO(), s3Conf, img.Path); err != nil {
			return err
		}
	}

	// after we are done with deleting images from s3 bucket delete main record
	return database.DB.Unscoped().Delete(&models.CarPostsModel{}, ID).Error
}
