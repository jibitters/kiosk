// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package database

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/test/containers"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

CREATE INDEX idx_tickets_owner_ticket_importance_level_ticket_status_updated_at ON tickets (owner, ticket_importance_level, ticket_status, updated_At);

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

CREATE INDEX idx_comments_ticket_id_created_at ON comments (ticket_id, created_at);`

func runPostgresContainer() (testcontainers.Container, int, error) {
	containerPort, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return nil, 0, err
	}

	request := testcontainers.ContainerRequest{
		Image:        "postgres:11",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_DB": "kiosk", "POSTGRES_USER": "kiosk", "POSTGRES_PASSWORD": "password"},
		WaitingFor:   wait.ForListeningPort(containerPort),
	}

	container, err := containers.NewContainer(request)
	if err != nil {
		return nil, 0, err
	}

	mappedPort, err := container.MappedPort(containers.ContainersContext, containerPort)
	if err != nil {
		return nil, 0, err
	}

	return container, mappedPort.Int(), nil
}

func TestBuildConnectionString(t *testing.T) {
	connectionString := buildConnectionString("localhost", 5432, "kiosk", "kiosk", "password", 10, "enable")

	if connectionString != "host=localhost port=5432 dbname=kiosk user=kiosk password=password connect_timeout=10 sslmode=enable" {
		t.Logf("Actual: %v Expected: host=localhost port=5432 dbname=kiosk user=kiosk password=password connect_timeout=10 sslmode=enable", connectionString)
		t.FailNow()
	}
}

func TestMigrate(t *testing.T) {
	container, port, err := runPostgresContainer()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)

	directory, err := ioutil.TempDir("", "migration")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	first, err := ioutil.TempFile(directory, "1_*.up.sql")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer first.Close()

	first.WriteString(firstMigrationSchema)

	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:               "localhost",
		Port:               port,
		Name:               "kiosk",
		User:               "kiosk",
		Password:           "password",
		ConnectionTimeout:  10,
		MaxConnection:      8,
		SSLMode:            "disable",
		MigrationDirectory: "file://" + filepath.Dir(first.Name()),
	}}

	if err := Migrate(config); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
}

func TestMigrate_ConnectionFailure(t *testing.T) {
	directory, err := ioutil.TempDir("", "migration")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	first, err := ioutil.TempFile(directory, "1_*.up.sql")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer first.Close()

	first.WriteString(firstMigrationSchema)

	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:               "localhost",
		Port:               5432 + 5432,
		Name:               "kiosk",
		User:               "kiosk",
		Password:           "password",
		ConnectionTimeout:  10,
		MaxConnection:      8,
		SSLMode:            "disable",
		MigrationDirectory: "file://" + filepath.Dir(first.Name()),
	}}

	if err := Migrate(config); err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}
}

func TestMigrate_SQLSyntaxError(t *testing.T) {
	container, port, err := runPostgresContainer()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)

	directory, err := ioutil.TempDir("", "migration")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	first, err := ioutil.TempFile(directory, "1_*.up.sql")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer first.Close()

	first.WriteString(`
	-- Tickets table definition.
	CRETE TABLE tickets (
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
	);`)

	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:               "localhost",
		Port:               port,
		Name:               "kiosk",
		User:               "kiosk",
		Password:           "password",
		ConnectionTimeout:  10,
		MaxConnection:      8,
		SSLMode:            "disable",
		MigrationDirectory: "file://" + filepath.Dir(first.Name()),
	}}

	if err := Migrate(config); err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}
}

func TestMigrate_NoChange(t *testing.T) {
	container, port, err := runPostgresContainer()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)

	directory, err := ioutil.TempDir("", "migration")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	first, err := ioutil.TempFile(directory, "1_*.up.sql")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer first.Close()

	first.WriteString(firstMigrationSchema)

	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:               "localhost",
		Port:               port,
		Name:               "kiosk",
		User:               "kiosk",
		Password:           "password",
		ConnectionTimeout:  10,
		MaxConnection:      8,
		SSLMode:            "disable",
		MigrationDirectory: "file://" + filepath.Dir(first.Name()),
	}}

	if err := Migrate(config); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	if err := Migrate(config); err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
}
