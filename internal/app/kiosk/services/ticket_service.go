// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package services

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4"
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
	request.TicketStatus = rpc.TicketStatus_NEW
	if err := service.validateCreate(request); err != nil {
		return nil, err
	}

	if err := service.insertOne(request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Read returns back a ticket by using the provided id to find one.
func (service *TicketService) Read(context context.Context, request *rpc.Id) (*rpc.Ticket, error) {
	if request.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "read_ticket.invalid_id")
	}

	ticket, err := service.findOne(request.Id)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// Update updates a ticket by using the provided values.
func (service *TicketService) Update(context context.Context, request *rpc.Ticket) (*empty.Empty, error) {
	if err := service.validateUpdate(request); err != nil {
		return nil, err
	}

	if err := service.updateOne(request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes a ticket by using the provided id. Use carefully, it has hard delete effect on db.
func (service *TicketService) Delete(context context.Context, request *rpc.Id) (*empty.Empty, error) {
	if request.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "delete_ticket.invalid_id")
	}

	if err := service.deleteOne(request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Filter returns a paginated response of tickets filtered by provided values.
func (service *TicketService) Filter(context context.Context, request *rpc.FilterTicketsRequest) (*rpc.FilterTicketsResponse, error) {
	// TODO: Complete implementation.

	return &rpc.FilterTicketsResponse{}, nil
}

func (service *TicketService) insertOne(ticket *rpc.Ticket) error {
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

func (service *TicketService) findOne(id int64) (*rpc.Ticket, error) {
	findTicketQuery := `SELECT * FROM tickets WHERE id = $1`
	findCommentsQuery := `SELECT * FROM comments WHERE comments.ticket_id = $1`

	batch := &pgx.Batch{}
	batch.Queue(findTicketQuery, id)
	batch.Queue(findCommentsQuery, id)

	results := service.db.SendBatch(context.Background(), batch)
	defer results.Close()

	ticket := &rpc.Ticket{}
	ticketImportanceLevel := ""
	ticketStatus := ""
	issuedAt := new(time.Time)
	updatedAt := new(time.Time)

	row := results.QueryRow()
	if err := row.Scan(
		&ticket.Id,
		&ticket.Issuer,
		&ticket.Owner,
		&ticket.Subject,
		&ticket.Content,
		&ticket.Metadata,
		&ticketImportanceLevel,
		&ticketStatus,
		issuedAt,
		updatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "read_ticket.not_found")
		}

		service.logger.Error("error on finding a ticket: %v", err)
		return nil, status.Error(codes.Internal, "read_ticket.failed")
	}

	ticket.TicketImportanceLevel = rpc.TicketImportanceLevel(rpc.TicketImportanceLevel_value[ticketImportanceLevel])
	ticket.TicketStatus = rpc.TicketStatus(rpc.TicketStatus_value[ticketStatus])
	ticket.IssuedAt = issuedAt.Format(time.RFC3339Nano)
	ticket.UpdatedAt = updatedAt.Format(time.RFC3339Nano)

	rows, err := results.Query()
	if err != nil {
		service.logger.Error("error on finding comments: %v", err)
		return nil, status.Error(codes.Internal, "read_ticket.failed")
	}
	defer rows.Close()

	for rows.Next() {
		id := int64(0)
		ticketID := int64(0)
		owner := ""
		content := ""
		metadata := ""
		createdAt := new(time.Time)
		updatedAt := new(time.Time)

		err := rows.Scan(&id, &ticketID, &owner, &content, &metadata, createdAt, updatedAt)
		if err != nil {
			service.logger.Error("error on scanning rows: %v", err)
			return nil, status.Error(codes.Internal, "read_ticket.failed")
		}

		ticket.Comments = append(ticket.Comments, &rpc.Comment{
			Id:        id,
			Owner:     owner,
			Content:   content,
			Metadata:  metadata,
			CreatedAt: createdAt.Format(time.RFC3339Nano),
			UpdatedAt: updatedAt.Format(time.RFC3339Nano),
		})
	}

	return ticket, nil
}

func (service *TicketService) updateOne(ticket *rpc.Ticket) error {
	query := `UPDATE tickets SET ticket_status=$1 WHERE id=$2`

	commandTag, err := service.db.Exec(context.Background(), query, ticket.TicketStatus.String(), ticket.Id)
	if err != nil {
		service.logger.Error("error on updating ticket: %v", err)
		return status.Error(codes.Internal, "update_ticket.failed")
	}

	if commandTag.RowsAffected() != 1 {
		return status.Error(codes.NotFound, "update_ticket.not_found")
	}

	return nil
}

func (service *TicketService) deleteOne(id *rpc.Id) error {
	deleteCommentsQuery := `DELETE FROM comments WHERE ticket_id=$1`
	deleteTicketQuery := `DELETE FROM tickets WHERE id=$1`

	batch := &pgx.Batch{}
	batch.Queue("BEGIN")
	batch.Queue(deleteCommentsQuery, id.Id)
	batch.Queue(deleteTicketQuery, id.Id)
	batch.Queue("COMMIT")

	results := service.db.SendBatch(context.Background(), batch)
	if err := results.Close(); err != nil {
		service.logger.Error("error on deleting ticket: %v", err)
		return status.Error(codes.Internal, "delete_ticket.failed")
	}

	return nil
}

func (service *TicketService) validateCreate(ticket *rpc.Ticket) error {
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

	if ticket.TicketStatus != rpc.TicketStatus_NEW {
		return status.Error(codes.InvalidArgument, "create_ticket.invalid_status")
	}

	return nil
}

func (service *TicketService) validateUpdate(ticket *rpc.Ticket) error {
	if ticket.Id < 1 {
		return status.Error(codes.InvalidArgument, "update_ticket.invalid_id")
	}

	if ticket.TicketStatus == rpc.TicketStatus_NEW {
		return status.Error(codes.InvalidArgument, "update_ticket.invalid_ticket_status")
	}

	return nil
}
