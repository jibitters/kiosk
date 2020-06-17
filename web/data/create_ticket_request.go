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
func (ctr *CreateTicketRequest) Validate() *errors.Type {
	if len(ctr.Issuer) == 0 {
		return errors.InvalidArgument("issuer.is_required", "")
	}

	if len(ctr.Owner) == 0 {
		return errors.InvalidArgument("owner.is_required", "")
	}

	if len(ctr.Subject) == 0 {
		return errors.InvalidArgument("subject.is_required", "")
	}

	if len(ctr.Content) == 0 {
		return errors.InvalidArgument("content.is_required", "")
	}

	if ctr.ImportanceLevel != models.TicketImportanceLevelLow &&
		ctr.ImportanceLevel != models.TicketImportanceLevelMedium &&
		ctr.ImportanceLevel != models.TicketImportanceLevelHigh &&
		ctr.ImportanceLevel != models.TicketImportanceLevelCritical {

		return errors.InvalidArgument("importanceLevel.not_valid", "")
	}

	return nil
}
