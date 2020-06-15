package models

import "time"

// Model is a basic database model abstraction that only includes required columns for all models.
type Model struct {
	ID         int64
	CreatedAt  time.Time
	ModifiedAt time.Time
}
