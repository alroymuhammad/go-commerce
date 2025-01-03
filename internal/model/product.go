package model

import (
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"` // Changed from int64 to float64
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}
