package models

import (
	"gorm.io/gorm"
)

type CarPostsModel struct {
	gorm.Model
	Year          int     `json:"year" gorm:"not null"`
	Description   string  `json:"description"`
	Mileage       int     `json:"mileage" gorm:"not null"`
	Price         float32 `json:"price" gorm:"not null"`
	ExteriorColor string  `json:"exterior_color" gorm:"not null"`
	InteriorColor string  `json:"interior_color" gorm:"not null"`
	Vin           string  `json:"vin" gorm:"not null"`
	Address       string  `json:"address" gorm:"not null"`
	SellerID      uint    `json:"seller_id" gorm:"not null"`
	CarModelID    uint    `json:"car_model_id" gorm:"not null"`
	CardPhotoURL  string  `json:"card_photo_url" gorm:"-"`

	CarModel   CarModelsModel   `json:"car_model" gorm:"foreignKey:CarModelID"`                                                // Many-To-One
	PostImages []CarImagesModel `json:"post_images" gorm:"foreignKey:CarPostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // One-To-Many
}

func (CarPostsModel) TableName() string {
	return "car_posts"
}
