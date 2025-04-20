package response

import "github.com/KaungHtetHein116/personal-task-manager/api/v1/models"

type ProjectsResponse struct {
	ID          uint
	Name        string
	Tasks       []models.Task
	Description *string
}
