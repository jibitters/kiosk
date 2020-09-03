package models

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

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
	q := `INSERT INTO tickets (issuer, owner, subject, content, metadata, importance_level, status, created_at,
			modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW());`

	_, e := r.db.Exec(ctx, q, ticket.Issuer, ticket.Owner, ticket.Subject, ticket.Content, ticket.Metadata,
		ticket.ImportanceLevel, TicketStatusNew)
	if e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return et
	}

	return nil
}

// LoadByID tries to load a ticket and its comments from tickets table.
func (r *TicketRepository) LoadByID(ctx context.Context, id int64) (*Ticket, *errors.Type) {
	q := `SELECT id, issuer, owner, subject, content, metadata, importance_level, status, created_at, modified_at
			FROM tickets WHERE id = $1;`

	commentsQ := `SELECT id, ticket_id, owner, content, metadata, created_at, modified_at FROM comments WHERE
					ticket_id = $1 ORDER BY created_at DESC;`

	batch := &pgx.Batch{}
	batch.Queue(q, id)
	batch.Queue(commentsQ, id)

	results := r.db.SendBatch(ctx, batch)
	defer func() { _ = results.Close() }()

	ticket := &Ticket{}
	var metadata sql.NullString

	row := results.QueryRow()
	e := row.Scan(&ticket.ID, &ticket.Issuer, &ticket.Owner, &ticket.Subject, &ticket.Content, &metadata,
		&ticket.ImportanceLevel, &ticket.Status, &ticket.CreatedAt, &ticket.ModifiedAt)
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

		e := rows.Scan(&comment.ID, &comment.TicketID, &comment.Owner, &comment.Content, &metadata, &comment.CreatedAt,
			&comment.ModifiedAt)
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
	q := `UPDATE tickets SET subject = $1, metadata = $2, importance_level = $3, status = $4, modified_at = NOW()
			WHERE id = $5;`

	command, e := r.db.Exec(ctx, q, ticket.Subject, ticket.Metadata, ticket.ImportanceLevel, ticket.Status, ticket.ID)
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

// Filter tries to filter tickets. If there is another page of result when loading tickets, the second returned value
// will be true, otherwise false.
func (r *TicketRepository) Filter(ctx context.Context, issuer, owner string, importanceLevel TicketImportanceLevel,
	status TicketStatus, fromDate, toDate string, pageNumber, pageSize int) ([]*Ticket, bool, *errors.Type) {

	q, args := r.buildFilterQuery(issuer, owner, importanceLevel, status, fromDate, toDate, pageNumber, pageSize)
	rows, e := r.db.Query(ctx, q, args...)
	if e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return nil, false, et
	}
	defer rows.Close()

	tickets := make([]*Ticket, 0)
	ticketsMap := make(map[int64]*Ticket)
	for rows.Next() {
		ticket := &Ticket{}
		var metadata sql.NullString

		e := rows.Scan(&ticket.ID, &ticket.Issuer, &ticket.Owner, &ticket.Subject, &ticket.Content, &metadata,
			&ticket.ImportanceLevel, &ticket.Status, &ticket.CreatedAt, &ticket.ModifiedAt)
		if e != nil {
			et := errors.InternalServerError("unknown", "")
			r.logger.Error(et.FingerPrint, ": ", e.Error())
			return nil, false, et
		}

		if metadata.Valid {
			ticket.Metadata = metadata.String
		}

		tickets = append(tickets, ticket)
		ticketsMap[ticket.ID] = ticket
	}

	hasNextPage := len(tickets) > pageSize
	if hasNextPage {
		// Drop the extra one.
		tickets = tickets[:len(tickets)-1]
	}

	if len(tickets) > 0 {
		q, args = r.buildLoadCommentsQuery(tickets)
		rows, e = r.db.Query(ctx, q, args...)
		if e != nil {
			et := errors.InternalServerError("unknown", "")
			r.logger.Error(et.FingerPrint, ": ", e.Error())
			return nil, false, et
		}
		defer rows.Close()

		for rows.Next() {
			comment := &Comment{}
			var metadata sql.NullString

			e := rows.Scan(&comment.ID, &comment.TicketID, &comment.Owner, &comment.Content, &metadata,
				&comment.CreatedAt, &comment.ModifiedAt)
			if e != nil {
				et := errors.InternalServerError("unknown", "")
				r.logger.Error(et.FingerPrint, ": ", e.Error())
				return nil, false, et
			}

			if metadata.Valid {
				comment.Metadata = metadata.String
			}

			ticketsMap[comment.TicketID].Comments = append(ticketsMap[comment.TicketID].Comments, comment)
		}
	}

	return tickets, hasNextPage, nil
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

func (r *TicketRepository) buildFilterQuery(issuer, owner string, importanceLevel TicketImportanceLevel,
	status TicketStatus, fromDate, toDate string, pageNumber, pageSize int) (string, []interface{}) {

	offset := (pageNumber - 1) * pageSize
	limit := pageSize

	args := make([]interface{}, 0)
	q := strings.Builder{}

	q.WriteString(`SELECT id, issuer, owner, subject, content, metadata, importance_level, status, created_at,
						modified_at FROM tickets WHERE`)

	counter := 0
	counter++
	q.WriteString(` modified_at >= $` + strconv.Itoa(counter))
	args = append(args, fromDate)

	counter++
	q.WriteString(` AND modified_at < $` + strconv.Itoa(counter))
	args = append(args, toDate)

	if issuer != "" {
		counter++
		q.WriteString(` AND issuer = $` + strconv.Itoa(counter))
		args = append(args, issuer)
	}

	if owner != "" {
		counter++
		q.WriteString(` AND owner = $` + strconv.Itoa(counter))
		args = append(args, owner)
	}

	if importanceLevel != "" {
		counter++
		q.WriteString(` AND importance_level = $` + strconv.Itoa(counter))
		args = append(args, importanceLevel)
	}

	if status != "" {
		counter++
		q.WriteString(` AND status = $` + strconv.Itoa(counter))
		args = append(args, status)
	}

	counter++
	q.WriteString(` ORDER BY modified_at DESC OFFSET $` + strconv.Itoa(counter))
	args = append(args, offset)

	counter++
	q.WriteString(` LIMIT $` + strconv.Itoa(counter))
	args = append(args, limit+1)

	return q.String(), args
}

func (r *TicketRepository) buildLoadCommentsQuery(tickets []*Ticket) (string, []interface{}) {
	q := strings.Builder{}
	args := make([]interface{}, 0)

	q.WriteString(`SELECT id, ticket_id, owner, content, metadata, created_at, modified_at FROM comments WHERE
						ticket_id IN (`)

	counter := 0
	for _, t := range tickets {
		if counter > 0 {
			q.WriteString(`, `)
		}
		counter++
		q.WriteString(`$`)
		q.WriteString(strconv.Itoa(counter))

		args = append(args, t.ID)
	}

	q.WriteString(`) ORDER BY created_at DESC;`)

	return q.String(), args
}
