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
	Comments        []*CommentResponse           `json:"comments"`
	CreatedAt       string                       `json:"createdAt"`
	ModifiedAt      string                       `json:"modifiedAt"`
}

// LoadFromTicket loads the fields of current model from provided ticket.
func (tr *TicketResponse) LoadFromTicket(ticket *models.Ticket) {
	tr.ID = ticket.ID
	tr.Issuer = ticket.Issuer
	tr.Owner = ticket.Owner
	tr.Subject = ticket.Subject
	tr.Content = ticket.Content
	tr.Metadata = ticket.Metadata
	tr.ImportanceLevel = ticket.ImportanceLevel
	tr.Status = ticket.Status

	for _, c := range ticket.Comments {
		cr := &CommentResponse{}
		cr.LoadFromComment(c)
		tr.Comments = append(tr.Comments, cr)
	}

	tr.CreatedAt = ticket.CreatedAt.Format(time.RFC3339Nano)
	tr.ModifiedAt = ticket.ModifiedAt.Format(time.RFC3339Nano)
}

// CommentResponse model definition.
type CommentResponse struct {
	ID         int64  `json:"ID"`
	TicketID   int64  `json:"ticketID"`
	Owner      string `json:"owner"`
	Content    string `json:"content"`
	Metadata   string `json:"metadata"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
}

// LoadFromComment loads the fields of current model from provided comment.
func (cr *CommentResponse) LoadFromComment(comment *models.Comment) {
	cr.ID = comment.ID
	cr.TicketID = comment.TicketID
	cr.Owner = comment.Owner
	cr.Content = comment.Content
	cr.Metadata = comment.Metadata
	cr.CreatedAt = comment.CreatedAt.Format(time.RFC3339Nano)
	cr.ModifiedAt = comment.ModifiedAt.Format(time.RFC3339Nano)
}
