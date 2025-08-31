package services

import (
	"car_sales/internal/models"
	"car_sales/internal/repositories"
)

/*
Services handle the bussiness logic.
They coordinate multiple repositories (check if user exists + save to DB)
Hide complexities from handler
*/

func CreateCarBrand(carBrand *models.CarBrandsModel) error {
	return repositories.CreateCarBrand(carBrand)
}

func GetAllCarBrands() ([]*models.CarBrandsModel, error) {
	return repositories.GetAllCarBrands()
}

func GetCarBrandById(id int) (*models.CarBrandsModel, error) {
	return repositories.GetCarBrandById(id)
}

func UpdateCarBrand(carBrand *models.CarBrandsModel) error {
	return repositories.UpdateCarBrand(carBrand)
}

func DeleteCarBrand(id int) error {
	return repositories.DeleteCarBrand(id)
}
