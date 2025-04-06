package models

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Description string
	Status      bool
	ProjectID   uint
	Project     Project `gorm:"foreignkey:ProjectID"`
}
