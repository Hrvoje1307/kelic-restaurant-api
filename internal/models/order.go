package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	MenuItemID int     `json:"menu_item_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

type Order struct {
	ID         int         `json:"id"`
	TableID    int         `json:"table_id"`
	Items      []OrderItem `json:"items"`
	Status     OrderStatus `json:"status"`
	TotalPrice float64     `json:"total_price"`
	Notes      string      `json:"notes,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
