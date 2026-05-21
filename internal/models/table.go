package models

import "time"

type Table struct {
	ID          string    `json:"id"`
	TableNumber int       `json:"table_number"`
	Capacity    int       `json:"capacity"`
	Location    string    `json:"location,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type TableInput struct {
	TableNumber int    `json:"table_number" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	Location    string `json:"location" binding:"omitempty,oneof=unutra terasa vip"`
	IsActive    *bool  `json:"is_active"`
}
