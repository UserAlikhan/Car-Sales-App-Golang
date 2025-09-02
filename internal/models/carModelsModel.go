package models

import "gorm.io/gorm"

type CarModelsModel struct {
	gorm.Model
	Name       string          `json:"name" gorm:"not null; unique"`
	CarBrandID uint            `json:"car_brand_id" gorm:"not null"`
	CarPosts   []CarPostsModel `json:"car_posts" gorm:"foreignKey:CarModelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CarModelsModel) TableName() string {
	return "car_models"
}
