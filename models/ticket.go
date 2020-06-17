package models

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/errors"
	"go.uber.org/zap"
)

// Ticket is the entity model of tickets table.
type Ticket struct {
	Model

	Issuer          string
	Owner           string
	Subject         string
	Content         string
	Metadata        string
	ImportanceLevel TicketImportanceLevel
	Status          TicketStatus
	Comments        []*Comment
}

// TicketRepository is the repository implementation of Ticket model.
type TicketRepository struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

// NewTicketRepository returns back a newly created and ready to use TicketRepository.
func NewTicketRepository(logger *zap.SugaredLogger, db *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{logger: logger, db: db}
}

// Insert tries to insert a ticket into tickets table.
func (r *TicketRepository) Insert(ctx context.Context, ticket Ticket) *errors.Type {
	q := `INSERT INTO tickets (
			issuer,
			owner,
			subject,
			content,
			metadata,
			importance_level,
			status,
			created_at,
			modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW());`

	_, e := r.db.Exec(ctx, q,
		ticket.Issuer,
		ticket.Owner,
		ticket.Subject,
		ticket.Content,
		ticket.Metadata,
		ticket.ImportanceLevel,
		TicketStatusNew,
	)

	if e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return et
	}

	return nil
}

// LoadByID tries to load a ticket from tickets table.
func (r *TicketRepository) LoadByID(ctx context.Context, id int64) (*Ticket, *errors.Type) {
	q := `SELECT
			id,
			issuer,
			owner,
			subject,
			content,
			metadata,
			importance_level,
			status,
			created_at,
			modified_at FROM tickets WHERE id = $1;`

	commentsQ := `SELECT
					id,
					ticket_id,
					owner,
					content,
					metadata,
					created_at,
					modified_at FROM comments WHERE ticket_id = $1 ORDER BY created_at DESC;`

	batch := &pgx.Batch{}
	batch.Queue(q, id)
	batch.Queue(commentsQ, id)

	results := r.db.SendBatch(ctx, batch)
	defer func() { _ = results.Close() }()

	ticket := &Ticket{}
	var metadata sql.NullString

	row := results.QueryRow()
	e := row.Scan(
		&ticket.ID,
		&ticket.Issuer,
		&ticket.Owner,
		&ticket.Subject,
		&ticket.Content,
		&metadata,
		&ticket.ImportanceLevel,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.ModifiedAt,
	)

	if e != nil {
		if e == pgx.ErrNoRows {
			return nil, errors.NotFound("ticket.not_found", "")
		}

		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return nil, et
	}

	if metadata.Valid {
		ticket.Metadata = metadata.String
	}

	rows, e := results.Query()
	if e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return nil, et
	}
	defer rows.Close()

	for rows.Next() {
		comment := &Comment{}
		var metadata sql.NullString

		e := rows.Scan(
			&comment.ID,
			&comment.TicketID,
			&comment.Owner,
			&comment.Content,
			&metadata,
			&comment.CreatedAt,
			&comment.ModifiedAt,
		)
		if e != nil {
			et := errors.InternalServerError("unknown", "")
			r.logger.Error(et.FingerPrint, ": ", e.Error())
			return nil, et
		}

		if metadata.Valid {
			comment.Metadata = metadata.String
		}

		ticket.Comments = append(ticket.Comments, comment)
	}

	return ticket, nil
}

// Update tries to update a ticket record.
func (r *TicketRepository) Update(ctx context.Context, ticket *Ticket) *errors.Type {
	q := `UPDATE tickets
			SET subject = $1,
				metadata = $2,
				importance_level = $3,
				status = $4,
				modified_at = NOW()
			WHERE id = $5;`

	command, e := r.db.Exec(ctx, q,
		ticket.Subject,
		ticket.Metadata,
		ticket.ImportanceLevel,
		ticket.Status,
		ticket.ID,
	)

	if e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return et
	}

	if command.RowsAffected() == 0 {
		et := errors.PreconditionFailed("ticket.not_found", "")
		return et
	}

	return nil
}

// DeleteByID tries to delete a ticket and all of its comments.
func (r *TicketRepository) DeleteByID(ctx context.Context, id int64) *errors.Type {
	begin := `BEGIN;`
	commentsQ := `DELETE FROM comments WHERE ticket_id=$1;`
	q := `DELETE FROM tickets WHERE id=$1;`
	commit := `COMMIT;`

	batch := &pgx.Batch{}
	batch.Queue(begin)
	batch.Queue(commentsQ, id)
	batch.Queue(q, id)
	batch.Queue(commit)

	results := r.db.SendBatch(ctx, batch)
	if e := results.Close(); e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return et
	}

	return nil
}

// TicketImportanceLevel model.
type TicketImportanceLevel string

// Different importance level instances.
const (
	TicketImportanceLevelLow      TicketImportanceLevel = "LOW"
	TicketImportanceLevelMedium   TicketImportanceLevel = "MEDIUM"
	TicketImportanceLevelHigh     TicketImportanceLevel = "HIGH"
	TicketImportanceLevelCritical TicketImportanceLevel = "CRITICAL"
)

// TicketStatus model.
type TicketStatus string

// Different ticket status instances.
const (
	TicketStatusNew      TicketStatus = "NEW"
	TicketStatusReplied  TicketStatus = "REPLIED"
	TicketStatusResolved TicketStatus = "RESOLVED"
	TicketStatusClosed   TicketStatus = "CLOSED"
	TicketStatusBlocked  TicketStatus = "BLOCKED"
)
