// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package services

import (
	"context"
	"fmt"
	"strings"

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
	newCommentSMSTemplate          = "A new comment created for ticket with ID: %d. Please check the panel."
	newCommentEmailSubjectTemplate = "A new comment created for ticket with ID: %d"
	newCommentEmailBodyTemplate    = ""
)

// CommentService is the implementation of comment related rpc methods.
type CommentService struct {
	config *configuration.Config
	logger *logging.Logger
	db     *pgxpool.Pool
	nats   *natsclient.Conn
}

// NewCommentService creates and returns a new ready to use comment service implementation.
func NewCommentService(config *configuration.Config, logger *logging.Logger, db *pgxpool.Pool, nats *natsclient.Conn) *CommentService {
	return &CommentService{
		config: config,
		logger: logger,
		db:     db,
		nats:   nats,
	}
}

// Create creates a new comment with provided values.
func (service *CommentService) Create(context context.Context, request *rpc.Comment) (*empty.Empty, error) {
	if err := service.validateCreate(request); err != nil {
		return nil, err
	}

	if err := service.insertOne(context, request); err != nil {
		return nil, err
	}

	service.notify(request)

	return &empty.Empty{}, nil
}

// Update updates a comment by using the provided values.
func (service *CommentService) Update(context context.Context, request *rpc.Comment) (*empty.Empty, error) {
	if err := service.updateOne(context, request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes a comment by using the provided id. Use carefully, it has hard delete effect on database.
func (service *CommentService) Delete(context context.Context, request *rpc.Id) (*empty.Empty, error) {
	if err := service.validateDelete(request); err != nil {
		return nil, err
	}

	if err := service.deleteOne(context, request); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (service *CommentService) validateCreate(request *rpc.Comment) error {
	request.Owner = strings.TrimSpace(request.Owner)
	request.Content = strings.TrimSpace(request.Content)
	request.Metadata = strings.TrimSpace(request.Metadata)

	if len(request.Owner) == 0 {
		return status.Error(codes.InvalidArgument, "create_comment.empty_owner")
	}

	if len(request.Content) == 0 {
		return status.Error(codes.InvalidArgument, "create_comment.empty_content")
	}

	return nil
}

func (service *CommentService) validateDelete(request *rpc.Id) error {
	if request.Id < 1 {
		return status.Error(codes.InvalidArgument, "delete_comment.invalid_id")
	}

	return nil
}

func (service *CommentService) insertOne(context context.Context, comment *rpc.Comment) error {
	insertCommentQuery := `INSERT INTO comments(ticket_id, owner, content, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, now(), now())`
	updateTicketQuery := `UPDATE tickets SET updated_at=now() WHERE id=$1`

	batch := &pgx.Batch{}
	batch.Queue("BEGIN")
	batch.Queue(insertCommentQuery, comment.TicketId, comment.Owner, comment.Content, comment.Metadata)
	batch.Queue(updateTicketQuery, comment.TicketId)
	batch.Queue("COMMIT")

	results := service.db.SendBatch(context, batch)
	if err := results.Close(); err != nil {
		if strings.Contains(err.Error(), "comments_ticket_id_fkey") {
			return status.Error(codes.InvalidArgument, "create_comment.ticket_not_exists")
		}

		service.logger.Error("error on inserting new comment: %v", err)
		return status.Error(codes.Internal, "create_comment.failed")
	}

	return nil
}

func (service *CommentService) updateOne(context context.Context, comment *rpc.Comment) error {
	query := `UPDATE comments SET metadata=$1, updated_at=now() WHERE id=$2`

	commandTag, err := service.db.Exec(context, query, comment.Metadata, comment.Id)
	if err != nil {
		service.logger.Error("error on updating comment: %v", err)
		return status.Error(codes.Internal, "update_comment.failed")
	}

	if commandTag.RowsAffected() != 1 {
		return status.Error(codes.NotFound, "update_comment.not_found")
	}

	return nil
}

func (service *CommentService) deleteOne(context context.Context, id *rpc.Id) error {
	query := `DELETE FROM comments WHERE id=$1`

	_, err := service.db.Exec(context, query, id.Id)
	if err != nil {
		service.logger.Error("error on deleting comment: %v", err)
		return status.Error(codes.Internal, "delete_comment.failed")
	}

	return nil
}

func (service *CommentService) notify(request *rpc.Comment) {
	protobytes, _ := proto.Marshal(&notifiermodels.NotificationRequest{
		NotificationType: notifiermodels.NotificationRequest_Type(notifiermodels.NotificationRequest_Type_value[service.config.Notifications.Comment.New.Type]),
		Message:          fmt.Sprintf(newCommentSMSTemplate, request.TicketId),
		Subject:          fmt.Sprintf(newCommentEmailSubjectTemplate, request.TicketId),
		Body:             newCommentEmailBodyTemplate,
		Sender:           service.config.Notifications.Comment.New.Sender,
		Cc:               service.config.Notifications.Comment.New.CC,
		Bcc:              service.config.Notifications.Comment.New.BCC,
		Recipient:        service.config.Notifications.Comment.New.Recipients,
	})

	_ = service.nats.Publish(service.config.Notifier.Subject, protobytes)
}
