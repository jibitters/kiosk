package database

import (
	"testing"

	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/test/containers"
)

func TestConnectToDatabase(t *testing.T) {
	container, port, err := startPostgresContainer()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)

	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:              "localhost",
		Port:              port,
		Name:              "kiosk",
		User:              "kiosk",
		Password:          "password",
		ConnectionTimeout: 10,
		MaxConnection:     8,
		SSLMode:           "disable",
	}}

	db, err := ConnectToDatabase(config)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer db.Close()
}

func TestConnectToDatabase_ParseError(t *testing.T) {
	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:              "localhost",
		Port:              -1,
		Name:              "kiosk",
		User:              "invalid",
		Password:          "invalid",
		ConnectionTimeout: 10,
		MaxConnection:     8,
		SSLMode:           "disable",
	}}

	db, err := ConnectToDatabase(config)
	if err == nil {
		db.Close()
		t.Logf("Expected error here!")
		t.FailNow()
	}
}

func TestConnectToDatabase_InvalidCredentials(t *testing.T) {
	container, port, err := startPostgresContainer()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)

	config := &configuration.Config{Postgres: configuration.PostgresConfig{
		Host:              "localhost",
		Port:              port,
		Name:              "kiosk",
		User:              "invalid",
		Password:          "invalid",
		ConnectionTimeout: 10,
		MaxConnection:     8,
		SSLMode:           "disable",
	}}

	db, err := ConnectToDatabase(config)
	if err == nil {
		db.Close()
		t.Logf("Expected error here!")
		t.FailNow()
	}
}
