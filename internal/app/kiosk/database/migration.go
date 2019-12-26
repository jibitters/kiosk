// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package database

import (
	"database/sql"
	"fmt"

	migration "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/pkg/logging"
)

// Migrate tries to connect to a postgres instance with the connection information provided for migration.
func Migrate(logger *logging.Logger, config *configuration.Config) error {
	connectionString := buildConnectionString(config.Postgres.Host, config.Postgres.Port, config.Postgres.Name, config.Postgres.User, config.Postgres.Password, config.Postgres.ConnectionTimeout, config.Postgres.SSLMode)

	db, err := openConnection(connectionString)
	if err != nil {
		return err
	}

	if err := pingConnection(db); err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error on closing database connection: %v", err)
		}
	}()

	if err := migrate(db, config.Postgres.MigrationDirectory); err != nil {
		return err
	}

	return nil
}

func openConnection(connection string) (*sql.DB, error) {
	return sql.Open("postgres", connection)
}

func pingConnection(db *sql.DB) error {
	return db.Ping()
}

func migrate(db *sql.DB, migrationDirectory string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrate, err := migration.NewWithDatabaseInstance(migrationDirectory, "postgres", driver)
	if err != nil {
		return err
	}

	if err := migrate.Up(); err != nil {
		if err == migration.ErrNoChange {
			return nil
		}

		return err
	}

	return nil
}

func buildConnectionString(host string, port int, name, user, password string, connectionTimeout int, sslMode string) string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s connect_timeout=%d sslmode=%s",
		host,
		port,
		name,
		user,
		password,
		connectionTimeout,
		sslMode,
	)
}
