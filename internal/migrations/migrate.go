package main

import (
	"car_sales/internal/configs"
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func init() {
	configs.LoadEnvVariables()
	database.ConnectToDB()
}

func main() {
	database.DB.AutoMigrate(
		&models.CarBrandsModel{},
		&models.CarModelsModel{},
		&models.UsersModel{},
		&models.CarPostsModel{},
		&models.CarImagesModel{},
	)
}
