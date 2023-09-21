package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

var UserRequest struct {
	Email    string
	Password string
}
