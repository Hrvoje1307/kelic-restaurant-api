package models

import "time"

type MenuCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

type MenuCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}
