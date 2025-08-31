package models

type CarModelsModel struct {
	Name       string `json:"name" gorm:"not null; unique"`
	CarBrandID uint   `json:"car_brand_id" gorm:"not null"`
}

func (CarModelsModel) TableName() string {
	return "car_models"
}
