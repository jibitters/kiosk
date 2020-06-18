package data

import (
	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
)

// UpdateCommentRequest model definition.
type UpdateCommentRequest struct {
	ID       int64  `json:"ID"`
	Metadata string `json:"metadata"`
}

// Validate validates the request.
func (r *UpdateCommentRequest) Validate() *errors.Type {

	return nil
}

// AsComment converts this request model into comment model.
func (r *UpdateCommentRequest) AsComment() *models.Comment {
	return &models.Comment{
		Model:    models.Model{ID: r.ID},
		Metadata: r.Metadata,
	}
}
