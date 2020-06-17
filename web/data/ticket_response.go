package data

import (
	"time"

	"github.com/jibitters/kiosk/models"
)

// TicketResponse model definition.
type TicketResponse struct {
	ID              int64                        `json:"ID"`
	Issuer          string                       `json:"issuer"`
	Owner           string                       `json:"owner"`
	Subject         string                       `json:"subject"`
	Content         string                       `json:"content"`
	Metadata        string                       `json:"metadata"`
	ImportanceLevel models.TicketImportanceLevel `json:"importanceLevel"`
	Status          models.TicketStatus          `json:"status"`
	Comments        []CommentResponse            `json:"comments"`
	CreatedAt       time.Time                    `json:"createdAt"`
	ModifiedAt      time.Time                    `json:"modifiedAt"`
}

// CommentResponse model definition.
type CommentResponse struct {
	ID         int64     `json:"ID"`
	TicketID   int64     `json:"ticketID"`
	Owner      string    `json:"owner"`
	Content    string    `json:"content"`
	Metadata   string    `json:"metadata"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
