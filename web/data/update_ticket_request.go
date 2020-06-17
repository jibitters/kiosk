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
func (utr *UpdateTicketRequest) Validate() *errors.Type {
	if len(utr.Subject) == 0 {
		return errors.InvalidArgument("subject.is_required", "")
	}

	if utr.ImportanceLevel != models.TicketImportanceLevelLow &&
		utr.ImportanceLevel != models.TicketImportanceLevelMedium &&
		utr.ImportanceLevel != models.TicketImportanceLevelHigh &&
		utr.ImportanceLevel != models.TicketImportanceLevelCritical {

		return errors.InvalidArgument("importanceLevel.not_valid", "")
	}

	if utr.Status != models.TicketStatusReplied &&
		utr.Status != models.TicketStatusResolved &&
		utr.Status != models.TicketStatusClosed &&
		utr.Status != models.TicketStatusBlocked {

		return errors.InvalidArgument("status.not_valid", "")
	}

	return nil
}

// AsTicket converts this request model into ticket model.
func (utr *UpdateTicketRequest) AsTicket() *models.Ticket {
	return &models.Ticket{
		Model:           models.Model{ID: utr.ID},
		Subject:         utr.Subject,
		Metadata:        utr.Metadata,
		ImportanceLevel: utr.ImportanceLevel,
		Status:          utr.Status,
	}
}
