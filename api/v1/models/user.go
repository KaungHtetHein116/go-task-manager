package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `gorm:"not null" json:"-"`
	Projects []Project
}
