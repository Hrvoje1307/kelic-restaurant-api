package models

type Category string

const (
	CategoryStarter   Category = "starter"
	CategoryMain      Category = "main"
	CategoryDessert   Category = "dessert"
	CategoryDrink     Category = "drink"
)

type MenuItem struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    Category `json:"category"`
	Available   bool     `json:"available"`
}
