package postgres

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lireza/lib/configuring"
	"go.uber.org/zap"
)

// Connect tries to connect to a postgres instance with the information provided in config instance.
func Connect(logger *zap.SugaredLogger, config *configuring.Config) (*pgxpool.Pool, error) {
	connectionString := config.Get("db.postgres.connection_string").
		StringOrElse("postgres://localhost:5432/kiosk?sslmode=disable")

	minPoolConnections := config.Get("db.postgres.pool_min_connections").
		IntOrElse(2)

	maxPoolConnections := config.Get("db.postgres.pool_max_connections").
		IntOrElse(8)

	migrationDirectory := config.Get("db.postgres.migration_directory").
		StringOrElse("file://migration/postgres")

	logger.Debug("db.postgres.connection_string -> ", connectionString)
	logger.Info("db.postgres.pool_min_connections -> ", minPoolConnections)
	logger.Info("db.postgres.pool_max_connections -> ", maxPoolConnections)
	logger.Info("db.postgres.migration_directory -> ", migrationDirectory)

	dbConfig, e := pgxpool.ParseConfig(connectionString)
	if e != nil {
		return nil, e
	}

	dbConfig.MinConns = int32(minPoolConnections)
	dbConfig.MaxConns = int32(maxPoolConnections)

	db, e := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if e != nil {
		return nil, e
	}

	return db, nil
}

// Migrate tries to connect to a postgres instance and then runs database migration.
func Migrate(logger *zap.SugaredLogger, config *configuring.Config) error {
	connectionString := config.Get("db.postgres.connection_string").
		StringOrElse("postgres://localhost:5432/kiosk?sslmode=disable")

	migrationDirectory := config.Get("db.postgres.migration_directory").
		StringOrElse("file://migration/postgres")

	migratory, e := migrate.New(migrationDirectory, connectionString)
	if e != nil {
		return e
	}

	if e := migratory.Up(); e != nil && e != migrate.ErrNoChange {
		return e
	}

	logger.Info("Successfully executed database migration.")
	return nil
}
