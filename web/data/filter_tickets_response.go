package data

// FilterTicketsResponse model definition.
type FilterTicketsResponse struct {
	Tickets     []TicketResponse `json:"tickets"`
	HasNextPage bool             `json:"hasNextPage"`
}
