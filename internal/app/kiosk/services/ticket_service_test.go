// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package services

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/app/kiosk/database"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	"github.com/jibitters/kiosk/test/containers"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const firstMigrationSchema = `
-- Tickets table definition.
	CREATE TABLE tickets (
	    id                                 BIGSERIAL NOT NULL,
	    issuer                             VARCHAR(40) NOT NULL,
	    owner                              VARCHAR(40) NOT NULL,
	    subject                            VARCHAR(255) NOT NULL,
	    content                            TEXT NOT NULL,
	    metadata                           TEXT,
	    ticket_importance_level            VARCHAR(20) NOT NULL,
	    ticket_status                      VARCHAR(20) NOT NULL,
	    issued_at                          TIMESTAMP NOT NULL,
	    updated_at                         TIMESTAMP NOT NULL,
	    PRIMARY KEY (id)
	);

	CREATE INDEX idx_tickets_issuer_issued_at ON tickets (issuer, issued_at DESC);
	CREATE INDEX idx_tickets_owner_issued_at ON tickets (owner, issued_at DESC);
	CREATE INDEX idx_tickets_ticket_importance_level_ticket_status ON tickets (ticket_importance_level, ticket_status);

	-- Comments table definition.
	CREATE TABLE comments (
	    id                                 BIGSERIAL NOT NULL,
	    ticket_id                          BIGINT REFERENCES tickets,
	    owner                              VARCHAR(40) NOT NULL,
	    content                            TEXT NOT NULL,
	    metadata                           TEXT,
	    created_at                         TIMESTAMP NOT NULL,
	    updated_at                         TIMESTAMP NOT NULL,
	    PRIMARY KEY (id)
	);

	CREATE INDEX idx_comments_ticket_id ON comments (ticket_id);
	CREATE INDEX idx_comments_owner_created_at ON comments (owner, created_at DESC);`

func setupPostgresAndRunMigration() (testcontainers.Container, *pgxpool.Pool, error) {
	// Starting postgres container.
	containerPort, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return nil, nil, err
	}

	request := testcontainers.ContainerRequest{
		Image:        "postgres:11",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_DB": "kiosk", "POSTGRES_USER": "kiosk", "POSTGRES_PASSWORD": "password"},
		WaitingFor:   wait.ForListeningPort(containerPort),
	}

	container, err := containers.NewContainer(request)
	if err != nil {
		return nil, nil, err
	}

	mappedPort, err := container.MappedPort(containers.ContainersContext, containerPort)
	if err != nil {
		return nil, nil, err
	}

	// Running database migration.
	directory, err := ioutil.TempDir("", "migration")
	if err != nil {
		return nil, nil, err
	}

	first, err := ioutil.TempFile(directory, "1_*.up.sql")
	if err != nil {
		return nil, nil, err
	}
	defer first.Close()

	first.WriteString(firstMigrationSchema)

	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:               "localhost",
		Port:               mappedPort.Int(),
		Name:               "kiosk",
		User:               "kiosk",
		Password:           "password",
		ConnectionTimeout:  10,
		MaxConnection:      8,
		SSLMode:            "disable",
		MigrationDirectory: "file://" + filepath.Dir(first.Name()),
	}}

	if err := database.Migrate(config); err != nil {
		return nil, nil, err
	}

	// Getting database connection pool.
	db, err := database.ConnectToDatabase(config)
	if err != nil {
		return nil, nil, err
	}

	return container, db, nil
}

func TestCreate_InvalidArgument(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	defer db.Close()

	service := NewTicketService(logging.New(logging.DebugLevel), db)

	ticket := &rpc.Ticket{
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}
	shouldReturnInvalidArgument(t, service, ticket, "create_ticket.empty_issuer")

	ticket = &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}
	shouldReturnInvalidArgument(t, service, ticket, "create_ticket.empty_owner")

	ticket = &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               " ",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}
	shouldReturnInvalidArgument(t, service, ticket, "create_ticket.empty_subject")

	ticket = &rpc.Ticket{
		Issuer:  "Jibit",
		Owner:   "09203091992",
		Subject: "Documentation",
		Content: "	",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}
	shouldReturnInvalidArgument(t, service, ticket, "create_ticket.empty_content")

	ticket = &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_RESOLVED,
	}
	shouldReturnInvalidArgument(t, service, ticket, "create_ticket.invalid_status")
}

func TestCreate_DatabaseConnectionFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	db.Close()

	service := NewTicketService(logging.New(logging.DebugLevel), db)

	ticket := &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}
	shouldReturnInternal(t, service, ticket, "create_ticket.failed")
}

func TestCreate_DatabaseNetworkFailure(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	containers.CloseContainer(container)
	defer db.Close()

	service := NewTicketService(logging.New(logging.DebugLevel), db)

	ticket := &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}
	shouldReturnInternal(t, service, ticket, "create_ticket.failed")
}

func TestCreate(t *testing.T) {
	container, db, err := setupPostgresAndRunMigration()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)
	defer db.Close()

	service := NewTicketService(logging.New(logging.DebugLevel), db)

	ticket := &rpc.Ticket{
		Issuer:                "Jibit",
		Owner:                 "09203091992",
		Subject:               "Documentation",
		Content:               "Hello, i need some help about your technical documentation.",
		Metadata:              "{\"owner_ip\": \"185.186.187.188\"}",
		TicketImportanceLevel: rpc.TicketImportanceLevel_HIGH,
		TicketStatus:          rpc.TicketStatus_NEW,
	}

	if _, err := service.Create(context.Background(), ticket); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
}

func shouldReturnInvalidArgument(t *testing.T, service *TicketService, ticket *rpc.Ticket, message string) {
	_, err := service.Create(context.Background(), ticket)
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

func shouldReturnInternal(t *testing.T, service *TicketService, ticket *rpc.Ticket, message string) {
	_, err := service.Create(context.Background(), ticket)
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
