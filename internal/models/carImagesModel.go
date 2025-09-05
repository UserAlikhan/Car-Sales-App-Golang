package models

import "gorm.io/gorm"

type CarImagesModel struct {
	gorm.Model
	Path      string `json:"path" gorm:"not null; unique"`
	CarPostID uint   `json:"car_post_id" gorm:"not null"`
}

func (CarImagesModel) TableName() string {
	return "car_images"
}
