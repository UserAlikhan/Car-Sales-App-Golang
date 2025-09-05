package models

import "gorm.io/gorm"

type CarImagesModel struct {
	gorm.Model
	Path      string `json:"path" gorm:"not null; unique"`
	CarPostId uint   `json:"car_post_id"`
}

func (CarImagesModel) TableName() string {
	return "car_images"
}
