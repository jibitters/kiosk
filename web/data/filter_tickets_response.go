package data

import (
	"github.com/jibitters/kiosk/models"
)

// FilterTicketsResponse model definition.
type FilterTicketsResponse struct {
	Tickets     []*TicketResponse `json:"tickets"`
	HasNextPage bool              `json:"hasNextPage"`
}

// LoadFromTickets loads the fields of current model from provided tickets.
func (ftr *FilterTicketsResponse) LoadFromTickets(tickets []*models.Ticket, HasNextPage bool) {
	for _, t := range tickets {
		ticketResponse := &TicketResponse{}
		ticketResponse.LoadFromTicket(t)
		ftr.Tickets = append(ftr.Tickets, ticketResponse)
	}

	ftr.HasNextPage = HasNextPage
}
