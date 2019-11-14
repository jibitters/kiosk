// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package services

import (
	"context"
	"database/sql"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TicketService is the implementation of ticket related rpc methods.
type TicketService struct {
	logger *logging.Logger
	db     *pgxpool.Pool
}

// NewTicketService creates and returns a new ready to use ticket service implementation.
func NewTicketService(logger *logging.Logger, db *pgxpool.Pool) *TicketService {
	return &TicketService{
		logger: logger,
		db:     db,
	}
}

// Create creates a new ticket with provided values.
func (service *TicketService) Create(context context.Context, request *rpc.Ticket) (*empty.Empty, error) {
	if err := validateCreate(request); err != nil {
		return nil, err
	}

	if err := service.insertTicket(request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Read returns back a ticket by using the provided id to find one.
func (service *TicketService) Read(context context.Context, request *rpc.Id) (*rpc.Ticket, error) {
	if request.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "read_ticket.invalid_id")
	}

	ticket, err := service.findTicket(request.Id)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// Update updates a ticket by using the provided values.
func (service *TicketService) Update(context context.Context, request *rpc.Ticket) (*empty.Empty, error) {
	if err := validateUpdate(request); err != nil {
		return nil, err
	}

	if err := service.updateTicket(request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes a ticket by using the provided id. Use carefully, it has hard delete effect on db.
func (service *TicketService) Delete(context context.Context, request *rpc.Id) (*empty.Empty, error) {
	if request.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "delete_ticket.invalid_id")
	}

	if err := service.deleteTicket(request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// List returns a paginated response of tickets that are related to an owner.
func (service *TicketService) List(context context.Context, request *rpc.ListTicketsRequest) (*rpc.ListTicketsResponse, error) {
	// TODO: Complete implementation.
	return &rpc.ListTicketsResponse{}, nil
}

func (service *TicketService) insertTicket(ticket *rpc.Ticket) error {
	query := `
	INSERT INTO tickets(issuer, owner, subject, content, metadata, ticket_importance_level, ticket_status, issued_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())`

	_, err := service.db.Exec(
		context.Background(),
		query,
		ticket.Issuer,
		ticket.Owner,
		ticket.Subject,
		ticket.Content,
		ticket.Metadata,
		ticket.TicketImportanceLevel.String(),
		ticket.TicketStatus.String(),
	)

	if err != nil {
		service.logger.Error("error on inserting new ticket: %v", err)
		return status.Error(codes.Internal, "create_ticket.failed")
	}

	return nil
}

func (service *TicketService) findTicket(id int64) (*rpc.Ticket, error) {
	query := `
	SELECT (id, issuer, owner, subject, content, metadata, ticket_importance_level, ticket_status, issued_at, updated_at)
	FROM tickets WHERE id = $1`

	ticket := &rpc.Ticket{}
	ticketImportanceLevel := ""
	ticketStatus := ""
	row := service.db.QueryRow(context.Background(), query, id)
	if err := row.Scan(
		&ticket.Id,
		&ticket.Issuer,
		&ticket.Owner,
		&ticket.Subject,
		&ticket.Content,
		&ticket.Metadata,
		&ticketImportanceLevel,
		&ticketStatus,
		&ticket.IssuedAt,
		&ticket.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "read_ticket.not_found")
		}

		service.logger.Error("error on finding a ticket: %v", err)
		return nil, status.Error(codes.Internal, "read_ticket.failed")
	}

	ticket.TicketImportanceLevel = rpc.TicketImportanceLevel(rpc.TicketImportanceLevel_value[ticketImportanceLevel])
	ticket.TicketStatus = rpc.TicketStatus(rpc.TicketStatus_value[ticketStatus])

	return ticket, nil
}

func (service *TicketService) updateTicket(ticket *rpc.Ticket) error {
	query := `UPDATE tickets SET ticket_status=$1 WHERE id=$2`

	commandTag, err := service.db.Exec(context.Background(), query, ticket.TicketStatus.String(), ticket.Id)
	if err != nil {
		service.logger.Error("error on updating ticket: %v", err)
		return status.Error(codes.Internal, "update_ticket.failed")
	}

	if commandTag.RowsAffected() == 0 {
		return status.Error(codes.NotFound, "update_ticket.not_found")
	}

	return nil
}

func (service *TicketService) deleteTicket(id *rpc.Id) error {
	query := `DELETE FROM tickets WHERE id=$1`

	commandTag, err := service.db.Exec(context.Background(), query, id)
	if err != nil {
		service.logger.Error("error on deleting ticket: %v", err)
		return status.Error(codes.Internal, "delete_ticket.failed")
	}

	if commandTag.RowsAffected() == 0 {
		return status.Error(codes.NotFound, "delete_ticket.not_found")
	}

	return nil
}

func validateCreate(ticket *rpc.Ticket) error {
	ticket.Issuer = strings.TrimSpace(ticket.Issuer)
	ticket.Owner = strings.TrimSpace(ticket.Owner)
	ticket.Subject = strings.TrimSpace(ticket.Subject)
	ticket.Content = strings.TrimSpace(ticket.Content)
	ticket.Metadata = strings.TrimSpace(ticket.Metadata)

	if len(ticket.Issuer) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_issuer")
	}

	if len(ticket.Owner) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_owner")
	}

	if len(ticket.Subject) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_subject")
	}

	if len(ticket.Content) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_content")
	}

	return nil
}

func validateUpdate(ticket *rpc.Ticket) error {
	if ticket.Id < 1 {
		return status.Error(codes.InvalidArgument, "update_ticket.invalid_id")
	}

	if ticket.TicketStatus == rpc.TicketStatus_NEW {
		return status.Error(codes.InvalidArgument, "update_ticket.invalid_ticket_status")
	}

	return nil
}
