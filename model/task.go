package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Description   *string    `json:"description,omitempty"`
	DueDate       time.Time  `json:"due_date"`
	Priority      string     `json:"priority"`
	Status        string     `json:"status"`
	AssignedTo    *uuid.UUID `json:"assigned_to,omitempty"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	CustomerID    *uuid.UUID `json:"customer_id,omitempty"`
	OpportunityID *uuid.UUID `json:"opportunity_id,omitempty"`
	ContactID     *uuid.UUID `json:"contact_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}
