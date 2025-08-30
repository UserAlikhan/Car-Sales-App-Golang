package models

import "gorm.io/gorm"

// GORM will create ID, created_at and updated_at automatically
type CarBrandsModel struct {
	gorm.Model
	Name      string `json:"name" gorm:"not null;unique"`
	LogoImage string `json:"logo_image"`
}

// specify table name explicitly, otherwise it will be cars_brands_models
func (CarBrandsModel) TableName() string {
	return "car_brands"
}
