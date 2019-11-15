-- Tickets table definition.
CREATE TABLE tickets (
    id                                 BIGSERIAL NOT NULL,
    issuer                             VARCHAR(40) NOT NULL,
    owner                              VARCHAR(40) NOT NULL,
    subject                            VARCHAR(100) NOT NULL,
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
CREATE INDEX idx_comments_owner_created_at ON comments (owner, created_at DESC);
