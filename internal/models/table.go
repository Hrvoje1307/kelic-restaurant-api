package models

type Table struct {
	ID          string `json:"id"`
	TableNumber int    `json:"table_number"`
	Capacity    int    `json:"capacity"`
	Location    string `json:"location,omitempty"`
	IsActive    bool   `json:"is_active"`
}
