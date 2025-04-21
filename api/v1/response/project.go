package response

import "github.com/KaungHtetHein116/personal-task-manager/internal/entity"

type ProjectsResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description *string       `json:"description,omitempty"`
	Tasks       []entity.Task `json:"tasks,omitempty"`
}

type ProjectDetailResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description *string       `json:"description,omitempty"`
	Tasks       []entity.Task `json:"tasks"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}
