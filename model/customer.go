package model

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID             uuid.UUID  `json:"id"`
	CompanyName    string     `json:"company_name"`
	Industry       *string    `json:"industry,omitempty"`
	Address        *string    `json:"address,omitempty"`
	City           *string    `json:"city,omitempty"`
	Province       *string    `json:"province,omitempty"`
	PostalCode     *string    `json:"postal_code,omitempty"`
	Phone          *string    `json:"phone,omitempty"`
	Website        *string    `json:"website,omitempty"`
	CustomerStatus string     `json:"customer_status"`
	CustomerType   *string    `json:"customer_type,omitempty"`
	AssignedTo     *uuid.UUID `json:"assigned_to,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	AnnualRevenue  *float64   `json:"annual_revenue,omitempty"`
	Tags           []string   `json:"tags,omitempty"`
}
