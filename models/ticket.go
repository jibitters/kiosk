package models

import (
	"github.com/jackc/pgx/v4/pgxpool"
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

// ImportanceLevel model.
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
