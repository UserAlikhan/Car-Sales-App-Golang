package models

import "gorm.io/gorm"

type UsersModel struct {
	gorm.Model
	Firstname      string `json:"first_name" gorm:"not null; size:100"`
	Lastname       string `json:"last_name" gorm:"not null; size:100"`
	Username       string `json:"username" gorm:"not null; unique"`
	Email          string `json:"email" gorm:"not null; unique"`
	PhoneNumber    string `json:"phone_number" gorm:"unique"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
	IsAdmin        bool   `json:"is_admin" gorm:"default: false"`
}

func (UsersModel) TableName() string {
	return "users"
}

type LoginDataModel struct {
	Username string
	Email    string
	Password string
}
