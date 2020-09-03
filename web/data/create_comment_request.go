package data

import (
	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
)

// CreateCommentRequest model definition.
type CreateCommentRequest struct {
	TicketID int64  `json:"ticketID"`
	Owner    string `json:"owner"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}

// Validate validates the request.
func (r *CreateCommentRequest) Validate() *errors.Type {
	if r.TicketID <= 0 {
		return errors.InvalidArgument("ticketID.invalid", "")
	}

	if len(r.Owner) == 0 {
		return errors.InvalidArgument("owner.is_required", "")
	}

	if len(r.Owner) > 50 {
		return errors.InvalidArgument("owner.invalid_length", "")
	}

	if len(r.Content) == 0 {
		return errors.InvalidArgument("content.is_required", "")
	}

	if len(r.Content) > 5000 {
		return errors.InvalidArgument("content.invalid_length", "")
	}

	return nil
}

// AsComment converts this request model into comment model.
func (r *CreateCommentRequest) AsComment() *models.Comment {
	return &models.Comment{
		TicketID: r.TicketID,
		Owner:    r.Owner,
		Content:  r.Content,
		Metadata: r.Metadata,
	}
}
