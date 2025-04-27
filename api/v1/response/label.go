package response

import (
	"time"

	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
)

// LabelResponse represents the response structure for a label
type LabelResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewLabelResponse creates a new LabelResponse from a Label entity
func NewLabelResponse(label *entity.Label) *LabelResponse {
	return &LabelResponse{
		ID:        label.ID,
		Name:      label.Name,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}
}

// NewLabelResponses creates a slice of LabelResponse from a slice of Label entities
func NewLabelResponses(labels []entity.Label) []LabelResponse {
	responses := make([]LabelResponse, len(labels))
	for i, label := range labels {
		responses[i] = *NewLabelResponse(&label)
	}
	return responses
}
