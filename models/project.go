package models

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description *string
	UserID      uint
	User        User `gorm:"foreignkey:UserID"`
	Tasks       []Task
}

type CreateProjectInput struct {
	Name        string `validate:"required,min=3"`
	Description string `validate:"omitempty,min=3"`
}

type ProjectsResponse struct {
	ID          uint
	Name        string
	Tasks       []Task
	Description *string
}
