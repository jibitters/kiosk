package data

import (
	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
)

// CreateTicketRequest model definition.
type CreateTicketRequest struct {
	Issuer          string                       `json:"issuer"`
	Owner           string                       `json:"owner"`
	Subject         string                       `json:"subject"`
	Content         string                       `json:"content"`
	Metadata        string                       `json:"metadata"`
	ImportanceLevel models.TicketImportanceLevel `json:"importanceLevel"`
}

// Validate validates the request.
func (r *CreateTicketRequest) Validate() *errors.Type {
	if len(r.Issuer) == 0 {
		return errors.InvalidArgument("issuer.is_required", "")
	}

	if len(r.Issuer) > 50 {
		return errors.InvalidArgument("issuer.invalid_length", "")
	}

	if len(r.Owner) == 0 {
		return errors.InvalidArgument("owner.is_required", "")
	}

	if len(r.Owner) > 50 {
		return errors.InvalidArgument("owner.invalid_length", "")
	}

	if len(r.Subject) == 0 {
		return errors.InvalidArgument("subject.is_required", "")
	}

	if len(r.Subject) > 255 {
		return errors.InvalidArgument("subject.invalid_length", "")
	}

	if len(r.Content) == 0 {
		return errors.InvalidArgument("content.is_required", "")
	}

	if len(r.Content) > 5000 {
		return errors.InvalidArgument("content.invalid_length", "")
	}

	if r.ImportanceLevel != models.TicketImportanceLevelLow &&
		r.ImportanceLevel != models.TicketImportanceLevelMedium &&
		r.ImportanceLevel != models.TicketImportanceLevelHigh &&
		r.ImportanceLevel != models.TicketImportanceLevelCritical {

		return errors.InvalidArgument("importanceLevel.not_valid", "")
	}

	return nil
}

// AsTicket converts this request model into ticket model.
func (r *CreateTicketRequest) AsTicket() *models.Ticket {
	return &models.Ticket{
		Issuer:          r.Issuer,
		Owner:           r.Owner,
		Subject:         r.Subject,
		Content:         r.Content,
		Metadata:        r.Metadata,
		ImportanceLevel: r.ImportanceLevel,
	}
}
