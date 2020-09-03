package models

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/errors"
	"go.uber.org/zap"
)

// Comment is the entity model of comments table.
type Comment struct {
	Model

	TicketID int64
	Owner    string
	Content  string
	Metadata string
}

// CommentRepository is the repository implementation of Comment model.
type CommentRepository struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

// NewCommentRepository returns back a newly created and ready to use CommentRepository.
func NewCommentRepository(logger *zap.SugaredLogger, db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{logger: logger, db: db}
}

// Insert tries to insert a comment into comments table.
func (r *CommentRepository) Insert(ctx context.Context, comment Comment) *errors.Type {
	q := `INSERT INTO comments (ticket_id, owner, content, metadata, created_at, modified_at) VALUES
			($1, $2, $3, $4, NOW(), NOW());`

	_, e := r.db.Exec(ctx, q, comment.TicketID, comment.Owner, comment.Content, comment.Metadata)
	if e != nil {
		if strings.Contains(e.Error(), "comments_ticket_id_fkey") {
			return errors.PreconditionFailed("ticket.not_exists", "")
		}

		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return et
	}

	return nil
}

// LoadByID tries to load a comment from comments table.
func (r *CommentRepository) LoadByID(ctx context.Context, id int64) (*Comment, *errors.Type) {
	q := `SELECT id, ticket_id, owner, content, metadata, created_at, modified_at FROM comments WHERE id = $1;`

	comment := &Comment{}
	var metadata sql.NullString

	row := r.db.QueryRow(ctx, q, id)
	e := row.Scan(&comment.ID, &comment.TicketID, &comment.Owner, &comment.Content, &metadata, &comment.CreatedAt,
		&comment.ModifiedAt)
	if e != nil {
		if e == pgx.ErrNoRows {
			return nil, errors.NotFound("comment.not_found", "")
		}

		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return nil, et
	}

	if metadata.Valid {
		comment.Metadata = metadata.String
	}

	return comment, nil
}

// Update tries to update a comment record.
func (r *CommentRepository) Update(ctx context.Context, comment *Comment) *errors.Type {
	q := `UPDATE comments SET metadata = $1, modified_at = NOW() WHERE id = $2;`

	command, e := r.db.Exec(ctx, q, comment.Metadata, comment.ID)
	if e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return et
	}

	if command.RowsAffected() == 0 {
		et := errors.NotFound("comment.not_found", "")
		return et
	}

	return nil
}

// DeleteByID tries to delete a comment from comments table.
func (r *CommentRepository) DeleteByID(ctx context.Context, id int64) *errors.Type {
	q := `DELETE FROM comments WHERE id=$1;`

	_, e := r.db.Exec(ctx, q, id)
	if e != nil {
		et := errors.InternalServerError("unknown", "")
		r.logger.Error(et.FingerPrint, ": ", e.Error())
		return et
	}

	return nil
}
