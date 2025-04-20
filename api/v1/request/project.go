package request

type CreateProjectInput struct {
	Name        string `validate:"required,min=3"`
	Description string `validate:"omitempty,min=3"`
}
