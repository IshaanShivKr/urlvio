package model

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID
	UserID		string
	URL         string
	Code        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccessCount int64
}
