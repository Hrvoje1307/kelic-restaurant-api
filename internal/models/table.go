package models

type TableStatus string

const (
	TableStatusFree     TableStatus = "free"
	TableStatusOccupied TableStatus = "occupied"
	TableStatusReserved TableStatus = "reserved"
)

type Table struct {
	ID       int         `json:"id"`
	Number   int         `json:"number"`
	Capacity int         `json:"capacity"`
	Status   TableStatus `json:"status"`
}
