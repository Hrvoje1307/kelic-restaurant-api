package models

import "time"

type MenuItem struct {
	ID          string    `json:"id"`
	CategoryID  *string   `json:"category_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    *string   `json:"image_url"`
	IsAvailable bool      `json:"is_available"`
	Allergens   []string  `json:"allergens"`
	Tags        []string  `json:"tags"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MenuItemInput struct {
	CategoryID  *string  `json:"category_id"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	ImageURL    string   `json:"image_url"`
	IsAvailable *bool    `json:"is_available"`
	Allergens   []string `json:"allergens"`
	Tags        []string `json:"tags"`
	SortOrder   int      `json:"sort_order"`
}
