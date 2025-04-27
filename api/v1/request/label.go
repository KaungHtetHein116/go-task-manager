package request

type CreateLabelInput struct {
	Name string `json:"name" validate:"required"`
}

type UpdateLabelInput struct {
	Name string `json:"name" validate:"required"`
}
