package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID           uuid.UUID       `json:"id"`
	UserID       *uuid.UUID      `json:"user_id,omitempty"`
	ActivityType string          `json:"activity_type"`
	EntityType   string          `json:"entity_type"`
	EntityID     uuid.UUID       `json:"entity_id"`
	Description  string          `json:"description"`
	CreatedAt    time.Time       `json:"created_at"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}
