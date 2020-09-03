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
	Metadata        string                       `json:"metadata,omitempty"`
	ImportanceLevel models.TicketImportanceLevel `json:"importanceLevel"`
	Status          models.TicketStatus          `json:"status"`
	Comments        []*CommentResponse           `json:"comments,omitempty"`
	CreatedAt       string                       `json:"createdAt"`
	ModifiedAt      string                       `json:"modifiedAt"`
}

// LoadFromTicket populates the fields of current model from provided ticket.
func (r *TicketResponse) LoadFromTicket(ticket *models.Ticket) {
	r.ID = ticket.ID
	r.Issuer = ticket.Issuer
	r.Owner = ticket.Owner
	r.Subject = ticket.Subject
	r.Content = ticket.Content
	r.Metadata = ticket.Metadata
	r.ImportanceLevel = ticket.ImportanceLevel
	r.Status = ticket.Status

	for _, c := range ticket.Comments {
		cr := &CommentResponse{}
		cr.LoadFromComment(c)
		r.Comments = append(r.Comments, cr)
	}

	r.CreatedAt = ticket.CreatedAt.Format(time.RFC3339Nano)
	r.ModifiedAt = ticket.ModifiedAt.Format(time.RFC3339Nano)
}

// CommentResponse model definition.
type CommentResponse struct {
	ID         int64  `json:"ID"`
	TicketID   int64  `json:"ticketID"`
	Owner      string `json:"owner"`
	Content    string `json:"content"`
	Metadata   string `json:"metadata,omitempty"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
}

// LoadFromComment populates the fields of current model from provided comment.
func (r *CommentResponse) LoadFromComment(comment *models.Comment) {
	r.ID = comment.ID
	r.TicketID = comment.TicketID
	r.Owner = comment.Owner
	r.Content = comment.Content
	r.Metadata = comment.Metadata
	r.CreatedAt = comment.CreatedAt.Format(time.RFC3339Nano)
	r.ModifiedAt = comment.ModifiedAt.Format(time.RFC3339Nano)
}
