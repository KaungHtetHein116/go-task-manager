package seed

import (
	"github.com/KaungHtetHein116/personal-task-manager/database"
	"github.com/KaungHtetHein116/personal-task-manager/models"
)

func Run() {
	db := database.DB

	// Create User
	user := models.User{
		Username: "kaung",
		Email:    "kaung@example.com",
		Password: "hashedpassword",
	}
	db.Create(&user)

	// Create Project
	project := models.Project{
		Name:   "Project Alpha",
		UserID: user.ID,
	}
	db.Create(&project)

	// Create Labels
	labelBug := models.Label{Name: "bug", Color: "red"}
	labelFeature := models.Label{Name: "feature", Color: "green"}
	db.Create(&labelBug)
	db.Create(&labelFeature)

	// Create Task with Labels
	task := models.Task{
		Title:       "Fix login issue",
		Description: "The login crashes on iOS",
		Status:      false,
		Priority:    "high",
		ProjectID:   project.ID,
		Labels:      []models.Label{labelBug, labelFeature},
		Project:     project,
	}
	db.Create(&task)
}
