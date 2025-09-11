package repositories

import (
	"car_sales/internal/database"
	"car_sales/internal/models"

	"gorm.io/gorm"
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
	// if creating the brand succeeds but creating the models fails, I will have
	// an orphan brand. A transaction avoids that.
	// So, it is either everything is created or nothing is created
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// separate carModels from carBrand
		carModels := append([]models.CarModelsModel(nil), carBrand.CarModels...)
		carBrand.CarModels = nil
		// create car brand first
		if err := tx.Create(carBrand).Error; err != nil {
			return err
		}
		// add CarBrandID to each car model
		for i := range carModels {
			carModels[i].CarBrandID = carBrand.ID
		}

		// create car models
		if len(carModels) > 0 {
			if err := tx.Create(&carModels).Error; err != nil {
				return err
			}
		}
		// store car models back to the car brand
		carBrand.CarModels = carModels

		return nil
	})
}

func SaveCarBrand(carBrand *models.CarBrandsModel) error {
	// .Save method creates new object if there is nothing with the same primary key
	// or if there is an object with the same id, it updates the record
	return database.DB.Save(carBrand).Error
}
