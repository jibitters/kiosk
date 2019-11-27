// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	notifiermodels "github.com/jibitters/kiosk/g/rpc/notifier"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	natsclient "github.com/nats-io/nats.go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	newTicketSMSTemplate          = "A new Ticket with %s Importance Level Created by %s. Please Check Panel."
	newTicketEmailSubjectTemplate = "A new Ticket with %s Importance Level Created by %s"
	newTicketEmailBodyTemplate    = ""
)

// TicketService is the implementation of ticket related rpc methods.
type TicketService struct {
	config *configuration.Config
	logger *logging.Logger
	db     *pgxpool.Pool
	nats   *natsclient.Conn
}

// NewTicketService creates and returns a new ready to use ticket service implementation.
func NewTicketService(config *configuration.Config, logger *logging.Logger, db *pgxpool.Pool, nats *natsclient.Conn) *TicketService {
	return &TicketService{
		config: config,
		logger: logger,
		db:     db,
		nats:   nats,
	}
}

// Create creates a new ticket with provided values.
func (service *TicketService) Create(context context.Context, request *rpc.Ticket) (*empty.Empty, error) {
	if err := service.validateCreate(request); err != nil {
		return nil, err
	}

	if err := service.insertOne(context, request); err != nil {
		return nil, err
	}

	service.notify(request)

	return &empty.Empty{}, nil
}

// Read returns back a ticket by using the provided id to find one.
func (service *TicketService) Read(context context.Context, request *rpc.Id) (*rpc.Ticket, error) {
	if err := service.validateRead(request); err != nil {
		return nil, err
	}

	ticket, err := service.findOne(context, request.Id)
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

	if err := service.updateOne(context, request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes a ticket by using the provided id. Use carefully, it has hard delete effect on database.
func (service *TicketService) Delete(context context.Context, request *rpc.Id) (*empty.Empty, error) {
	if err := service.validateDelete(request); err != nil {
		return nil, err
	}

	if err := service.deleteOne(context, request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Filter returns a paginated response of tickets filtered by provided values.
func (service *TicketService) Filter(context context.Context, request *rpc.FilterTicketsRequest) (*rpc.FilterTicketsResponse, error) {
	if err := service.validateFilter(request); err != nil {
		return nil, err
	}

	filterTickets, err := service.filter(context, request)
	if err != nil {
		return nil, err
	}

	return filterTickets, nil
}

func (service *TicketService) validateCreate(request *rpc.Ticket) error {
	request.Issuer = strings.TrimSpace(request.Issuer)
	request.Owner = strings.TrimSpace(request.Owner)
	request.Subject = strings.TrimSpace(request.Subject)
	request.Content = strings.TrimSpace(request.Content)
	request.Metadata = strings.TrimSpace(request.Metadata)

	if len(request.Issuer) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_issuer")
	}

	if len(request.Owner) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_owner")
	}

	if len(request.Subject) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_subject")
	}

	if len(request.Content) == 0 {
		return status.Error(codes.InvalidArgument, "create_ticket.empty_content")
	}

	if request.TicketStatus != rpc.TicketStatus_NEW {
		return status.Error(codes.InvalidArgument, "create_ticket.invalid_ticket_status")
	}

	return nil
}

func (service *TicketService) validateRead(request *rpc.Id) error {
	if request.Id < 1 {
		return status.Error(codes.InvalidArgument, "read_ticket.invalid_id")
	}

	return nil
}

func (service *TicketService) validateUpdate(request *rpc.Ticket) error {
	if request.Id < 1 {
		return status.Error(codes.InvalidArgument, "update_ticket.invalid_id")
	}

	if request.TicketStatus == rpc.TicketStatus_NEW {
		return status.Error(codes.InvalidArgument, "update_ticket.invalid_ticket_status")
	}

	return nil
}

func (service *TicketService) validateDelete(request *rpc.Id) error {
	if request.Id < 1 {
		return status.Error(codes.InvalidArgument, "delete_ticket.invalid_id")
	}

	return nil
}

func (service *TicketService) validateFilter(request *rpc.FilterTicketsRequest) error {
	if request.FromDate == "" {
		request.FromDate = "2000-01-01T00:00:00Z"
	}

	if request.ToDate == "" {
		request.ToDate = time.Now().UTC().Format(time.RFC3339Nano)
	}

	if request.Page.Number < 1 {
		return status.Error(codes.InvalidArgument, "filter_tickets.invalid_page_number")
	}

	if request.Page.Size < 1 || request.Page.Size > 200 {
		return status.Error(codes.InvalidArgument, "filter_tickets.invalid_page_size")
	}

	return nil
}

func (service *TicketService) insertOne(context context.Context, ticket *rpc.Ticket) error {
	query := `
	INSERT INTO tickets(issuer, owner, subject, content, metadata, ticket_importance_level, ticket_status, issued_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())`

	_, err := service.db.Exec(
		context,
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

func (service *TicketService) findOne(context context.Context, id int64) (*rpc.Ticket, error) {
	findTicketQuery := `SELECT id, issuer, owner, subject, content, metadata, ticket_importance_level, ticket_status, issued_at, updated_at FROM tickets WHERE id = $1`
	findCommentsQuery := `SELECT id, ticket_id, owner, content, metadata, created_at, updated_at FROM comments WHERE comments.ticket_id = $1`

	batch := &pgx.Batch{}
	batch.Queue(findTicketQuery, id)
	batch.Queue(findCommentsQuery, id)

	results := service.db.SendBatch(context, batch)
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
			TicketId:  ticketID,
			Owner:     owner,
			Content:   content,
			Metadata:  metadata,
			CreatedAt: createdAt.Format(time.RFC3339Nano),
			UpdatedAt: updatedAt.Format(time.RFC3339Nano),
		})
	}

	return ticket, nil
}

func (service *TicketService) updateOne(context context.Context, ticket *rpc.Ticket) error {
	query := `UPDATE tickets SET ticket_status=$1, updated_at=now() WHERE id=$2`

	commandTag, err := service.db.Exec(context, query, ticket.TicketStatus.String(), ticket.Id)
	if err != nil {
		service.logger.Error("error on updating ticket: %v", err)
		return status.Error(codes.Internal, "update_ticket.failed")
	}

	if commandTag.RowsAffected() != 1 {
		return status.Error(codes.NotFound, "update_ticket.not_found")
	}

	return nil
}

func (service *TicketService) deleteOne(context context.Context, id *rpc.Id) error {
	deleteCommentsQuery := `DELETE FROM comments WHERE ticket_id=$1`
	deleteTicketQuery := `DELETE FROM tickets WHERE id=$1`

	batch := &pgx.Batch{}
	batch.Queue("BEGIN")
	batch.Queue(deleteCommentsQuery, id.Id)
	batch.Queue(deleteTicketQuery, id.Id)
	batch.Queue("COMMIT")

	results := service.db.SendBatch(context, batch)
	if err := results.Close(); err != nil {
		service.logger.Error("error on deleting ticket: %v", err)
		return status.Error(codes.Internal, "delete_ticket.failed")
	}

	return nil
}

func (service *TicketService) filter(context context.Context, request *rpc.FilterTicketsRequest) (*rpc.FilterTicketsResponse, error) {
	response := &rpc.FilterTicketsResponse{Page: &rpc.Page{}}

	query, args := buildFilterTicketsQuery(request)
	rows, err := service.db.Query(context, query, args...)
	if err != nil {
		service.logger.Error("error on filtering tickets: %v", err)
		return nil, status.Error(codes.Internal, "filter_tickets.failed")
	}
	defer rows.Close()

	tickets := make([]*rpc.Ticket, 0)
	ticketsMap := make(map[int64]*rpc.Ticket)
	for rows.Next() {
		ticket := &rpc.Ticket{}
		ticketImportanceLevel := ""
		ticketStatus := ""
		issuedAt := new(time.Time)
		updatedAt := new(time.Time)

		err := rows.Scan(
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
		)
		if err != nil {
			service.logger.Error("error on scanning rows: %v", err)
			return nil, status.Error(codes.Internal, "filter_tickets.failed")
		}

		tickets = append(tickets, ticket)
		ticketsMap[ticket.Id] = ticket
	}

	response.Page.Number = request.Page.Number
	response.Page.Size = request.Page.Size
	response.Page.HasNext = len(tickets) > int(request.Page.Size)

	// Drop the extra one.
	if response.Page.HasNext {
		tickets = tickets[:len(tickets)-1]
	}

	if len(tickets) > 0 {
		query, args = buildFindAllCommentsQuery(tickets)
		rows, err = service.db.Query(context, query, args...)
		if err != nil {
			service.logger.Error("error on reading comments: %v", err)
			return nil, status.Error(codes.Internal, "filter_tickets.failed")
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
				return nil, status.Error(codes.Internal, "filter_tickets.failed")
			}

			ticketsMap[ticketID].Comments = append(ticketsMap[ticketID].Comments, &rpc.Comment{
				Id:        id,
				TicketId:  ticketID,
				Owner:     owner,
				Content:   content,
				Metadata:  metadata,
				CreatedAt: createdAt.Format(time.RFC3339Nano),
				UpdatedAt: updatedAt.Format(time.RFC3339Nano),
			})
		}
	}

	response.Tickets = tickets
	return response, nil
}

func buildFilterTicketsQuery(request *rpc.FilterTicketsRequest) (string, []interface{}) {
	offset := (request.Page.Number - 1) * request.Page.Size
	limit := request.Page.Size

	args := make([]interface{}, 0)
	query := strings.Builder{}
	query.WriteString(`SELECT id, issuer, owner, subject, content, metadata, ticket_importance_level, ticket_status, issued_at, updated_at FROM tickets WHERE`)

	counter := 0
	counter++
	query.WriteString(` ticket_importance_level = $` + strconv.Itoa(counter))
	args = append(args, request.TicketImportanceLevel.String())

	counter++
	query.WriteString(` AND ticket_status = $` + strconv.Itoa(counter))
	args = append(args, request.TicketStatus.String())

	counter++
	query.WriteString(` AND issued_at >= $` + strconv.Itoa(counter))
	args = append(args, request.FromDate)

	counter++
	query.WriteString(` AND issued_at < $` + strconv.Itoa(counter))
	args = append(args, request.ToDate)

	if request.Issuer != "" {
		counter++
		query.WriteString(` AND issuer = $` + strconv.Itoa(counter))
		args = append(args, request.Issuer)
	}

	if request.Owner != "" {
		counter++
		query.WriteString(` AND owner = $` + strconv.Itoa(counter))
		args = append(args, request.Owner)
	}

	counter++
	query.WriteString(` ORDER BY issued_at DESC OFFSET $` + strconv.Itoa(counter))
	args = append(args, offset)

	counter++
	query.WriteString(` LIMIT $` + strconv.Itoa(counter))
	args = append(args, limit+1)

	return query.String(), args
}

func buildFindAllCommentsQuery(tickets []*rpc.Ticket) (string, []interface{}) {
	args := make([]interface{}, 0)
	query := strings.Builder{}
	query.WriteString(`SELECT id, ticket_id, owner, content, metadata, created_at, updated_at FROM comments WHERE ticket_id IN (`)

	counter := 0
	for _, t := range tickets {
		if counter > 0 {
			query.WriteString(`, `)
		}
		counter++
		query.WriteString(`$`)
		query.WriteString(strconv.Itoa(counter))

		args = append(args, t.Id)
	}
	query.WriteString(`) ORDER BY created_at DESC`)

	return query.String(), args
}

func (service *TicketService) notify(request *rpc.Ticket) {
	switch request.TicketStatus {
	case rpc.TicketStatus_NEW:
		switch request.TicketImportanceLevel {
		case rpc.TicketImportanceLevel_LOW:
			protobytes, _ := proto.Marshal(&notifiermodels.NotificationRequest{
				Type:      notifiermodels.NotificationRequest_Type(notifiermodels.NotificationRequest_Type_value[service.config.Notifications.Ticket.New.Low.Type]),
				Message:   fmt.Sprintf(newTicketSMSTemplate, request.TicketImportanceLevel, request.Owner),
				Subject:   fmt.Sprintf(newTicketEmailSubjectTemplate, request.TicketImportanceLevel, request.Owner),
				Body:      newTicketEmailBodyTemplate,
				Recipient: service.config.Notifications.Ticket.New.Low.Recipients,
			})
			service.nats.Publish(service.config.Notifier.Subject, protobytes)

		case rpc.TicketImportanceLevel_MEDIUM:
			protobytes, _ := proto.Marshal(&notifiermodels.NotificationRequest{
				Type:      notifiermodels.NotificationRequest_Type(notifiermodels.NotificationRequest_Type_value[service.config.Notifications.Ticket.New.Medium.Type]),
				Message:   fmt.Sprintf(newTicketSMSTemplate, request.TicketImportanceLevel, request.Owner),
				Subject:   fmt.Sprintf(newTicketEmailSubjectTemplate, request.TicketImportanceLevel, request.Owner),
				Body:      newTicketEmailBodyTemplate,
				Recipient: service.config.Notifications.Ticket.New.Medium.Recipients,
			})
			service.nats.Publish(service.config.Notifier.Subject, protobytes)

		case rpc.TicketImportanceLevel_HIGH:
			protobytes, _ := proto.Marshal(&notifiermodels.NotificationRequest{
				Type:      notifiermodels.NotificationRequest_Type(notifiermodels.NotificationRequest_Type_value[service.config.Notifications.Ticket.New.High.Type]),
				Message:   fmt.Sprintf(newTicketSMSTemplate, request.TicketImportanceLevel, request.Owner),
				Subject:   fmt.Sprintf(newTicketEmailSubjectTemplate, request.TicketImportanceLevel, request.Owner),
				Body:      newTicketEmailBodyTemplate,
				Recipient: service.config.Notifications.Ticket.New.High.Recipients,
			})
			service.nats.Publish(service.config.Notifier.Subject, protobytes)

		case rpc.TicketImportanceLevel_CRITICAL:
			protobytes, _ := proto.Marshal(&notifiermodels.NotificationRequest{
				Type:      notifiermodels.NotificationRequest_Type(notifiermodels.NotificationRequest_Type_value[service.config.Notifications.Ticket.New.Critical.Type]),
				Message:   fmt.Sprintf(newTicketSMSTemplate, request.TicketImportanceLevel, request.Owner),
				Subject:   fmt.Sprintf(newTicketEmailSubjectTemplate, request.TicketImportanceLevel, request.Owner),
				Body:      newTicketEmailBodyTemplate,
				Recipient: service.config.Notifications.Ticket.New.Critical.Recipients,
			})
			service.nats.Publish(service.config.Notifier.Subject, protobytes)

		default:
			service.logger.Warning("no notifier handler for %s importance level!", request.TicketImportanceLevel.String())
		}
	default:
		service.logger.Warning("no notifier handler for %s status!", request.TicketStatus.String())
	}
}
