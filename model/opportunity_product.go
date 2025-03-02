package model

import (
	"time"

	"github.com/google/uuid"
)

type OpportunityProduct struct {
	ID            uuid.UUID `json:"id"`
	OpportunityID uuid.UUID `json:"opportunity_id"`
	ProductID     uuid.UUID `json:"product_id"`
	Quantity      int       `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	Discount      float64   `json:"discount"`
	Total         float64   `json:"total"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
