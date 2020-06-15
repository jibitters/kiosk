package models

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// Comment is the entity model of comments table.
type Comment struct {
	Model

	TicketID     int64
	Owner           string
	Content         string
	Metadata        string
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
