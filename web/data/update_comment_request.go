package data

import (
	"github.com/jibitters/kiosk/errors"
)

// UpdateCommentRequest model definition.
type UpdateCommentRequest struct {
	ID       int64  `json:"ID"`
	Metadata string `json:"metadata"`
}

// Validate validates the request.
func (ucr *UpdateCommentRequest) Validate() *errors.Type {

	return nil
}
