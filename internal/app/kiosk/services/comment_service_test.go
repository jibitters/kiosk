// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package services

import (
	"context"
	"testing"

	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	"github.com/jibitters/kiosk/test/containers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateComment_InvalidArgument(t *testing.T) {
	service := NewCommentService(logging.New(logging.DebugLevel), nil)

	comment := &rpc.Comment{
		TicketId: 1,
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}
	createCommentShouldReturnInvalidArgument(t, service, comment, "create_comment.empty_owner")

	comment = &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}
	createCommentShouldReturnInvalidArgument(t, service, comment, "create_comment.empty_content")
}

func TestCreateComment_TicketNotExists(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	defer db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}
	createCommentShouldReturnInvalidArgument(t, service, comment, "create_comment.ticket_not_exists")
}

func TestCreateComment_DatabaseConnectionFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}
	createCommentShouldReturnInternal(t, service, comment, "create_comment.failed")
}

func TestCreateComment_DatabaseNetworkFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}
	createCommentShouldReturnInternal(t, service, comment, "create_comment.failed")
}

func TestCreateComment(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	defer db.Close()

	ticketService := NewTicketService(logging.New(logging.DebugLevel), db)

	ticket := &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}

	if _, err := ticketService.Create(context.Background(), ticket); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}

	if _, err := service.Create(context.Background(), comment); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
}

func TestUpdateComment_Notfound(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	defer db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		Id:       1000,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "",
	}
	updateCommentShouldReturnNotfound(t, service, comment, "update_comment.not_found")
}

func TestUpdateComment_DatabaseConnectionFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}
	updateCommentShouldReturnInternal(t, service, comment, "update_comment.failed")
}

func TestUpdateComment_DatabaseNetworkFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	containers.CloseContainer(container)
	defer db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}
	updateCommentShouldReturnInternal(t, service, comment, "update_comment.failed")
}

func TestUpdateComment(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	defer db.Close()

	ticketService := NewTicketService(logging.New(logging.DebugLevel), db)

	ticket := &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}

	if _, err := ticketService.Create(context.Background(), ticket); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}

	if _, err := service.Create(context.Background(), comment); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	foundTicket, err := ticketService.findOne(context.Background(), 1)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	comment = &rpc.Comment{
		Id:       foundTicket.Comments[0].Id,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "",
	}

	if _, err := service.Update(context.Background(), comment); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	foundTicket, err = ticketService.findOne(context.Background(), 1)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	if foundTicket.Comments[0].Metadata != comment.Metadata {
		t.Logf("Actual: %v Expected: %v", foundTicket.Comments[0].Metadata, comment.Metadata)
		t.FailNow()
	}
}

func TestDeleteComment_InvalidArgument(t *testing.T) {
	service := NewCommentService(logging.New(logging.DebugLevel), nil)

	id := &rpc.Id{Id: 0}
	deleteCommentShouldReturnInvalidArgument(t, service, id, "delete_comment.invalid_id")
}

func TestDeleteComment_DatabaseConnectionFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	id := &rpc.Id{Id: 1}
	deleteCommentShouldReturnInternal(t, service, id, "delete_comment.failed")
}

func TestDeleteComment_DatabaseNetworkFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	containers.CloseContainer(container)
	defer db.Close()

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	id := &rpc.Id{Id: 1}
	deleteCommentShouldReturnInternal(t, service, id, "delete_comment.failed")
}

func TestDeleteComment(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	defer db.Close()

	ticketService := NewTicketService(logging.New(logging.DebugLevel), db)

	ticket := &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}

	if _, err := ticketService.Create(context.Background(), ticket); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	service := NewCommentService(logging.New(logging.DebugLevel), db)

	comment := &rpc.Comment{
		TicketId: 1,
		Owner:    "support@jibit.ir",
		Content:  "Hello, please find API related docs on website.",
		Metadata: "{\"owner_ip\": \"185.186.187.188\"}",
	}

	if _, err := service.Create(context.Background(), comment); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	foundTicket, err := ticketService.findOne(context.Background(), 1)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	id := &rpc.Id{Id: foundTicket.Comments[0].Id}
	if _, err := service.Delete(context.Background(), id); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	foundTicket, err = ticketService.findOne(context.Background(), 1)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	if len(foundTicket.Comments) != 0 {
		t.Logf("Actual: %v Expected: %v", len(foundTicket.Comments), 0)
		t.FailNow()
	}
}

func createCommentShouldReturnInvalidArgument(t *testing.T, service *CommentService, comment *rpc.Comment, message string) {
	_, err := service.Create(context.Background(), comment)
	if err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}

	status, ok := status.FromError(err)
	if !ok {
		t.Logf("The returned error is not compatible with gRPC error types.")
		t.FailNow()
	}

	if status.Code() != codes.InvalidArgument {
		t.Logf("Actual: %v Expected: %v", status.Code(), codes.InvalidArgument)
		t.FailNow()
	}

	if status.Message() != message {
		t.Logf("Actual: %v Expected: %v", status.Message(), message)
		t.FailNow()
	}
}

func createCommentShouldReturnInternal(t *testing.T, service *CommentService, comment *rpc.Comment, message string) {
	_, err := service.Create(context.Background(), comment)
	if err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}

	status, ok := status.FromError(err)
	if !ok {
		t.Logf("The returned error is not compatible with gRPC error types.")
		t.FailNow()
	}

	if status.Code() != codes.Internal {
		t.Logf("Actual: %v Expected: %v", status.Code(), codes.InvalidArgument)
		t.FailNow()
	}

	if status.Message() != message {
		t.Logf("Actual: %v Expected: %v", status.Message(), message)
		t.FailNow()
	}
}

func updateCommentShouldReturnNotfound(t *testing.T, service *CommentService, comment *rpc.Comment, message string) {
	_, err := service.Update(context.Background(), comment)
	if err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}

	status, ok := status.FromError(err)
	if !ok {
		t.Logf("The returned error is not compatible with gRPC error types.")
		t.FailNow()
	}

	if status.Code() != codes.NotFound {
		t.Logf("Actual: %v Expected: %v", status.Code(), codes.NotFound)
		t.FailNow()
	}

	if status.Message() != message {
		t.Logf("Actual: %v Expected: %v", status.Message(), message)
		t.FailNow()
	}
}

func updateCommentShouldReturnInternal(t *testing.T, service *CommentService, comment *rpc.Comment, message string) {
	_, err := service.Update(context.Background(), comment)
	if err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}

	status, ok := status.FromError(err)
	if !ok {
		t.Logf("The returned error is not compatible with gRPC error types.")
		t.FailNow()
	}

	if status.Code() != codes.Internal {
		t.Logf("Actual: %v Expected: %v", status.Code(), codes.Internal)
		t.FailNow()
	}

	if status.Message() != message {
		t.Logf("Actual: %v Expected: %v", status.Message(), message)
		t.FailNow()
	}
}

func deleteCommentShouldReturnInvalidArgument(t *testing.T, service *CommentService, id *rpc.Id, message string) {
	_, err := service.Delete(context.Background(), id)
	if err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}

	status, ok := status.FromError(err)
	if !ok {
		t.Logf("The returned error is not compatible with gRPC error types.")
		t.FailNow()
	}

	if status.Code() != codes.InvalidArgument {
		t.Logf("Actual: %v Expected: %v", status.Code(), codes.InvalidArgument)
		t.FailNow()
	}

	if status.Message() != message {
		t.Logf("Actual: %v Expected: %v", status.Message(), message)
		t.FailNow()
	}
}

func deleteCommentShouldReturnInternal(t *testing.T, service *CommentService, id *rpc.Id, message string) {
	_, err := service.Delete(context.Background(), id)
	if err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}

	status, ok := status.FromError(err)
	if !ok {
		t.Logf("The returned error is not compatible with gRPC error types.")
		t.FailNow()
	}

	if status.Code() != codes.Internal {
		t.Logf("Actual: %v Expected: %v", status.Code(), codes.Internal)
		t.FailNow()
	}

	if status.Message() != message {
		t.Logf("Actual: %v Expected: %v", status.Message(), message)
		t.FailNow()
	}
}
