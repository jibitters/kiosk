package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/db/postgres"
	"github.com/lireza/lib/configuring"
	"go.uber.org/zap"
)

// ConnectToDatabase connects to a postgres instance listening on provided host and port and then runs migration.
func ConnectToDatabase(host string, port int) (*pgxpool.Pool, error) {
	config := configuring.New()

	directory, e := ioutil.TempDir("", "migration")
	if e != nil {
		return nil, e
	}

	file, e := ioutil.TempFile(directory, "1_*.up.sql")
	if e != nil {
		return nil, e
	}

	defer func() { _ = file.Close() }()

	_, _ = file.WriteString(first)
	cs := fmt.Sprintf("postgres://user:password@%v:%v/kiosk?sslmode=disable", host, port)
	_ = os.Setenv("DB_POSTGRES_CONNECTION_STRING", cs)
	_ = os.Setenv("DB_POSTGRES_MIGRATION_DIRECTORY", "file://"+filepath.Dir(file.Name()))

	if e := postgres.Migrate(zap.S(), config); e != nil {
		return nil, e
	}

	db, e := postgres.Connect(zap.S(), config)
	if e != nil {
		return nil, e
	}

	return db, nil
}

var first = `
-- Tickets table definition.
CREATE TABLE tickets
(
    id               BIGSERIAL    NOT NULL,
    issuer           VARCHAR(50)  NOT NULL,
    owner            VARCHAR(50)  NOT NULL,
    subject          VARCHAR(255) NOT NULL,
    content          TEXT         NOT NULL,
    metadata         TEXT,
    importance_level VARCHAR(25)  NOT NULL,
    status           VARCHAR(25)  NOT NULL,
    created_at       TIMESTAMP    NOT NULL,
    modified_at      TIMESTAMP    NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX tickets_owner_importance_level_status_modified_at ON tickets (owner, importance_level, status, modified_at);

-- Comments table definition.
CREATE TABLE comments
(
    id          BIGSERIAL   NOT NULL,
    ticket_id   BIGINT REFERENCES tickets,
    owner       VARCHAR(50) NOT NULL,
    content     TEXT        NOT NULL,
    metadata    TEXT,
    created_at  TIMESTAMP   NOT NULL,
    modified_at TIMESTAMP   NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX comments_ticket_id_created_at ON comments (ticket_id, created_at);
`
