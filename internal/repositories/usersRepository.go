package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateUser(userData *models.UsersModel) (*models.UsersModel, error) {
	result := database.DB.Create(userData)
	if result.Error != nil {
		return nil, result.Error
	}

	return userData, nil
}

func GetUserById(id int) (*models.UsersModel, error) {
	var user *models.UsersModel
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByUsername(username string) (*models.UsersModel, error) {
	var user *models.UsersModel
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*models.UsersModel, error) {
	var user *models.UsersModel
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
