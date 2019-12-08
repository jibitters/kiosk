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

CREATE INDEX idx_comments_ticket_id_created_at ON comments (ticket_id, created_at);
