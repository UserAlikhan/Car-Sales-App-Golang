package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarBrand(carBrand *models.CarBrandsModel) (*models.CarBrandsModel, error) {
	result := database.DB.Create(carBrand)
	if result.Error != nil {
		return nil, result.Error
	}

	return carBrand, nil
}
