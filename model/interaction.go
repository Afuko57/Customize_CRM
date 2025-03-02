package model

import (
	"time"

	"github.com/google/uuid"
)

type Interaction struct {
	ID              uuid.UUID  `json:"id"`
	CustomerID      uuid.UUID  `json:"customer_id"`
	ContactID       *uuid.UUID `json:"contact_id,omitempty"`
	UserID          uuid.UUID  `json:"user_id"`
	InteractionType string     `json:"interaction_type"`
	Subject         string     `json:"subject"`
	Description     *string    `json:"description,omitempty"`
	InteractionDate time.Time  `json:"interaction_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	OpportunityID   *uuid.UUID `json:"opportunity_id,omitempty"`
	FollowUpDate    *time.Time `json:"follow_up_date,omitempty"`
	Status          *string    `json:"status,omitempty"`
}
