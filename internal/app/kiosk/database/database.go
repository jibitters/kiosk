// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
)

// ConnectToDatabase setups a connection to postgres instance.
func ConnectToDatabase(config *configuration.Config) (*pgxpool.Pool, error) {
	user := config.Postgres.User
	password := config.Postgres.Password
	host := config.Postgres.Host
	port := config.Postgres.Port
	name := config.Postgres.Name

	pgConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, host, port, name))
	if err != nil {
		return nil, err
	}

	pgConfig.MaxConns = int32(config.Postgres.MaxConnection)

	return pgxpool.ConnectConfig(context.Background(), pgConfig)
}
