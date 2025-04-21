package request

type CreateProjectInput struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateProjectInput struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"omitempty"`
}
