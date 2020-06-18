package data

import (
	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
)

// UpdateTicketRequest model definition.
type UpdateTicketRequest struct {
	ID              int64                        `json:"ID"`
	Subject         string                       `json:"subject"`
	Metadata        string                       `json:"metadata"`
	ImportanceLevel models.TicketImportanceLevel `json:"importanceLevel"`
	Status          models.TicketStatus          `json:"status"`
}

// Validate validates the request.
func (r *UpdateTicketRequest) Validate() *errors.Type {
	if len(r.Subject) == 0 {
		return errors.InvalidArgument("subject.is_required", "")
	}

	if r.ImportanceLevel != models.TicketImportanceLevelLow &&
		r.ImportanceLevel != models.TicketImportanceLevelMedium &&
		r.ImportanceLevel != models.TicketImportanceLevelHigh &&
		r.ImportanceLevel != models.TicketImportanceLevelCritical {

		return errors.InvalidArgument("importanceLevel.not_valid", "")
	}

	if r.Status != models.TicketStatusReplied &&
		r.Status != models.TicketStatusResolved &&
		r.Status != models.TicketStatusClosed &&
		r.Status != models.TicketStatusBlocked {

		return errors.InvalidArgument("status.not_valid", "")
	}

	return nil
}

// AsTicket converts this request model into ticket model.
func (r *UpdateTicketRequest) AsTicket() *models.Ticket {
	return &models.Ticket{
		Model:           models.Model{ID: r.ID},
		Subject:         r.Subject,
		Metadata:        r.Metadata,
		ImportanceLevel: r.ImportanceLevel,
		Status:          r.Status,
	}
}
