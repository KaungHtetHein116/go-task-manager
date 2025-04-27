package entity

import "github.com/jinzhu/gorm"

type Label struct {
	gorm.Model
	Name   string `gorm:"not null;unique"`
	UserID uint
	User   User   `gorm:"foreignkey:UserID"`
	Tasks  []Task `gorm:"many2many:task_label;"`
}
