package services

import (
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"car_sales/internal/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(userData *models.UsersModel) (*models.UsersModel, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userData.Password = string(hashedPassword)

	return repositories.CreateUser(userData)
}

func LoginUser(loginData *models.LoginDataModel) (string, error) {
	var user *models.UsersModel
	var err error

	// check if user with exists
	if loginData.Username != "" {
		user, err = GetUserByUsername(loginData.Username)
		if err != nil {
			return "", err
		}
	} else if loginData.Email != "" {
		user, err = GetUserByEmail(loginData.Email)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("Invalid credentials!")
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return "", err
	}

	tokenString, err := utils.GenerateToken(
		int(user.ID), loginData.Username,
		loginData.Email, user.IsAdmin,
	)
	if err != nil {
		return "", fmt.Errorf("Invalid to create a token.")
	}

	return tokenString, nil
}

func GetUserById(id int) (*models.UsersModel, error) {
	return repositories.GetUserById(id)
}

func GetUserByUsername(username string) (*models.UsersModel, error) {
	return repositories.GetUserByUsername(username)
}

func GetUserByEmail(email string) (*models.UsersModel, error) {
	return repositories.GetUserByEmail(email)
}
