package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarPost(carPost *models.CarPostsModel) (*models.CarPostsModel, error) {
	if err := database.DB.Create(&carPost).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("CarModel").First(carPost, carPost.ID).Error; err != nil {
		return nil, err
	}

	return carPost, nil
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

	result := database.DB.Preload("CarModel").Preload("PostImages").First(&carPost, ID)
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

func CountCarPostsTotalRows() int {
	var totalRows int64

	if err := database.DB.Model(&models.CarPostsModel{}).Count(&totalRows).Error; err != nil {
		return 0
	}

	return int(totalRows)
}

func GetCarPostsWithPagination(limit int, offset int) ([]*models.CarPostsModel, error) {
	var carPosts []*models.CarPostsModel

	// get the number of carPosts equal to the limit and
	// pass equal to the offset
	result := database.DB.
		Preload("CarModel").
		Preload("PostImages").
		Offset(offset).
		Limit(limit).
		Find(&carPosts)

	if result.Error != nil {
		return nil, result.Error
	}

	return carPosts, nil
}
