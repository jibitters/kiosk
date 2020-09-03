package data

import "github.com/jibitters/kiosk/models"

// FilterTicketsResponse model definition.
type FilterTicketsResponse struct {
	Tickets     []*TicketResponse `json:"tickets,omitempty"`
	HasNextPage bool              `json:"hasNextPage"`
}

// LoadFromTickets populates the fields of current model from provided tickets.
func (r *FilterTicketsResponse) LoadFromTickets(tickets []*models.Ticket, HasNextPage bool) {
	for _, t := range tickets {
		ticketResponse := &TicketResponse{}
		ticketResponse.LoadFromTicket(t)
		r.Tickets = append(r.Tickets, ticketResponse)
	}

	r.HasNextPage = HasNextPage
}
