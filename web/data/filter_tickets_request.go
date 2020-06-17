package data

import (
	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
)

// FilterTicketsRequest model definition.
type FilterTicketsRequest struct {
	Issuer          string                       `json:"issuer"`
	Owner           string                       `json:"owner"`
	ImportanceLevel models.TicketImportanceLevel `json:"importanceLevel"`
	Status          models.TicketStatus          `json:"status"`
	FromDate        string                       `json:"fromDate"`
	ToDate          string                       `json:"toDate"`
	PageNumber      int                          `json:"pageNumber"`
	PageSize        int                          `json:"pageSize"`
}

// Validate validates the request.
func (ftr *FilterTicketsRequest) Validate() *errors.Type {
	if len(ftr.Issuer) == 0 {
		return errors.InvalidArgument("issuer.is_required", "")
	}

	if ftr.PageNumber < 1 {
		return errors.InvalidArgument("pageNumber.not_valid", "")
	}

	if ftr.PageSize < 1 {
		return errors.InvalidArgument("pageSize.not_valid", "")
	}

	return nil
}
