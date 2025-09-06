package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarPost(carPost *models.CarPostsModel) (*models.CarPostsModel, error) {
	return carPost, database.DB.Create(&carPost).Error
}

func GetAllUsersCarPosts(userId uint) ([]*models.CarPostsModel, error) {
	var carPosts []*models.CarPostsModel

	result := database.DB.
		Where(&models.CarPostsModel{SellerID: userId}).
		Preload("CarModel").
		Preload("PostImages").
		Find(&carPosts)
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

func GetCarPostByIdWithPostImages(ID uint) (*models.CarPostsModel, error) {
	var carPost models.CarPostsModel

	result := database.DB.
		Preload("CarModel").
		Preload("CarModel.CarBrand").
		Preload("PostImages").
		First(&carPost, ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &carPost, nil
}

func DeleteCarPost(ID uint) error {
	// simply delete car post from database, car images records
	// will be delete right afterwards, because we have
	// CASCADE setted up
	return database.DB.Unscoped().Delete(&models.CarPostsModel{}, ID).Error
}
