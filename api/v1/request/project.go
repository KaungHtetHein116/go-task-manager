package request

type CreateProjectInput struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"omitempty,min=3"`
}

type UpdateProjectInput struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"omitempty,min=3"`
}
