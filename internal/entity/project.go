package entity

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description *string
	UserID      uint
	User        User `gorm:"foreignkey:UserID"`
	Tasks       []Task
}
