package models

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model
	Name   string `gorm:"unique;not null"`
	UserID uint
	User   User `gorm:"foreignkey:UserID"`
	Tasks  []Task
}

type CreateProjectInput struct {
	Name string `validate:"required,min=3"`
}
