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

	result := database.DB.Preload("CarModels").Find(&carBrands)
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
	carBrand, err := GetCarBrandById(id)
	if err != nil {
		return err
	}
	// .Unscoped is used to delete the records from db fully
	// without it GORM will just set a date for onDeleted row
	result := database.DB.Unscoped().Delete(&carBrand, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateCarBrandWithModels(carBrand *models.CarBrandsModel) error {
	// separate carModels from carBrand
	carModels := append([]models.CarModelsModel(nil), carBrand.CarModels...)
	carBrand.CarModels = nil
	// create car brand first
	if err := database.DB.Create(carBrand).Error; err != nil {
		return err
	}
	// add CarBrandID to each car model
	for i := range carModels {
		carModels[i].CarBrandID = carBrand.ID
	}

	// create car models
	if len(carModels) > 0 {
		if err := database.DB.Create(&carModels).Error; err != nil {
			return err
		}
	}

	carBrand.CarModels = carModels

	return nil
}
