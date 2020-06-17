package data

import (
	"github.com/jibitters/kiosk/errors"
)

// CreateCommentRequest model definition.
type CreateCommentRequest struct {
	TicketID int64  `json:"ticketID"`
	Owner    string `json:"owner"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}

// Validate validates the request.
func (ccr *CreateCommentRequest) Validate() *errors.Type {
	if len(ccr.Owner) == 0 {
		return errors.InvalidArgument("owner.is_required", "")
	}

	if len(ccr.Content) == 0 {
		return errors.InvalidArgument("content.is_required", "")
	}

	return nil
}
