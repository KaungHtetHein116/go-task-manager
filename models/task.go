package models

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Status      bool
	ProjectID   uint
	Priority    string  `sql:"type:ENUM('high', 'medium', 'low')"`
	Project     Project `gorm:"foreignkey:ProjectID"`
	Labels      []Label `gorm:"many2many:task_label;"`
}
