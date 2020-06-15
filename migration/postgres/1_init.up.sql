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
