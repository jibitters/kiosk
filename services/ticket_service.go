package services

import (
	"context"

	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/web/data"
)

// TicketService is the interface definition of ticket related functionalities.
type TicketService interface {
	Create(ctx context.Context, request *data.CreateTicketRequest) *errors.Type
	Load(ctx context.Context, ticketID *data.ID) (*data.TicketResponse, *errors.Type)
	Update(ctx context.Context, request *data.UpdateTicketRequest) *errors.Type
	Delete(ctx context.Context, ticketID *data.ID) *errors.Type
	Filter(ctx context.Context, request *data.FilterTicketsRequest) (*data.FilterTicketsResponse, *errors.Type)
}
