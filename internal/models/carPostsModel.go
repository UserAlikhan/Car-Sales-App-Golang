package models

import "gorm.io/gorm"

type CarPostsModel struct {
	gorm.Model
	Year          int     `json:"year"`
	Description   string  `json:"description"`
	Mileage       int     `json:"mileage"`
	Price         float32 `json:"price"`
	ExteriorColor string  `json:"exterior_color"`
	InteriorColor string  `json:"interior_color"`
	Vin           string  `json:"vin"`
	Address       string  `json:"address"`
	SellerID      uint    `json:"seller_id"`
	CarModelID    uint    `json:"car_model_id"`
}

func (CarPostsModel) TableName() string {
	return "car_posts"
}
