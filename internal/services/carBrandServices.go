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

func CreateCarBrand(carBrand *models.CarBrandsModel) (*models.CarBrandsModel, error) {
	return repositories.CreateCarBrand(carBrand)
}
