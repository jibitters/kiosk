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
    PRIMARY KEY  (id)
);

-- Comments table definition.
CREATE TABLE comments (
    id                                 BIGSERIAL NOT NULL,
    owner                              VARCHAR(40) NOT NULL,
    content                            TEXT NOT NULL,
    metadata                           TEXT,
    issued_at                          TIMESTAMP NOT NULL,
    updated_at                         TIMESTAMP NOT NULL,
    PRIMARY KEY  (id)
);

-- Ticket and its related comments table definition.
CREATE TABLE ticket_comments (
    ticket_id                          REFERENCES tickets,
    comment_id                         REFERENCES comments,
    PRIMARY KEY  (ticket_id, comment_id)
);
