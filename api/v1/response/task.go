package response

type TaskResponse struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      bool     `json:"status"`
	Priority    string   `json:"priority"`
	ProjectID   uint     `json:"project_id"`
	Labels      []string `json:"labels"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
