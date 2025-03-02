package model

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Position   *string   `json:"position,omitempty"`
	Email      *string   `json:"email,omitempty"`
	Phone      *string   `json:"phone,omitempty"`
	Mobile     *string   `json:"mobile,omitempty"`
	IsPrimary  bool      `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Notes      *string   `json:"notes,omitempty"`
}
