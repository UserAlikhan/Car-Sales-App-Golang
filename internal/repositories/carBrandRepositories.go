package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"
)

func CreateCarBrand(carBrand *models.CarBrandsModel) error {
	result := database.DB.Create(carBrand)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllCarBrands() ([]*models.CarBrandsModel, error) {
	var carBrands []*models.CarBrandsModel

	result := database.DB.Find(&carBrands)
	if result.Error != nil {
		return nil, result.Error
	}

	return carBrands, nil
}

func GetCarBrandById(id int) (*models.CarBrandsModel, error) {
	var carBrand *models.CarBrandsModel

	result := database.DB.First(&carBrand, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return carBrand, nil
}

func UpdateCarBrand(carBrand *models.CarBrandsModel) error {
	// check if car brand exists first
	_, err := GetCarBrandById(int(carBrand.ID))
	if err != nil {
		return err
	}

	result := database.DB.Model(&models.CarBrandsModel{}).Where("id = ?", carBrand.ID).Updates(carBrand)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteCarBrand(id int) error {
	// check if car brand exists first
	_, err := GetCarBrandById(id)
	if err != nil {
		return err
	}

	result := database.DB.Delete(&models.CarBrandsModel{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
