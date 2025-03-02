package model

import (
	"time"

	"github.com/google/uuid"
)

type Opportunity struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	CustomerID        uuid.UUID  `json:"customer_id"`
	ContactID         *uuid.UUID `json:"contact_id,omitempty"`
	Amount            *float64   `json:"amount,omitempty"`
	Stage             string     `json:"stage"`
	Probability       *int       `json:"probability,omitempty"`
	ExpectedCloseDate *time.Time `json:"expected_close_date,omitempty"`
	AssignedTo        *uuid.UUID `json:"assigned_to,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	Source            *string    `json:"source,omitempty"`
	Description       *string    `json:"description,omitempty"`
	Status            string     `json:"status"`
	ReasonLost        *string    `json:"reason_lost,omitempty"`
}
