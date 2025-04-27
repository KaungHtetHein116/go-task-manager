package request

type CreateTaskInput struct {
	Title       string   `validate:"required,min=3"`
	Description string   `validate:"omitempty,min=3"`
	Status      bool     `validate:"required"`
	ProjectID   uint     `json:"project_id" validate:"required"`
	UserID      uint     `json:"user_id"`
	Priority    string   `validate:"required,oneof=high medium low"`
	Labels      []string `json:"labels" validate:"omitempty"`
}

type UpdateTaskInput struct {
	Title       string   `json:"title" validate:"required,min=1,max=255"`
	Description string   `json:"description"`
	Status      bool     `json:"status" validate:"required"`
	Priority    string   `json:"priority" validate:"required,oneof=high medium low"`
	Labels      []string `json:"labels" validate:"omitempty"`
	UserID      uint     `json:"-"`
}
