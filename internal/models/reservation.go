package models

import "time"

type Reservation struct {
	ID           string    `json:"id"`
	TableID      *string   `json:"table_id"`
	GuestName    string    `json:"guest_name"`
	GuestEmail   string    `json:"guest_email"`
	GuestPhone   string    `json:"guest_phone,omitempty"`
	PartySize    int       `json:"party_size"`
	ReservedAt   time.Time `json:"reserved_at"`
	DurationMins int       `json:"duration_mins"`
	Status       string    `json:"status"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReservationInput struct {
	TableID      *string   `json:"table_id"`
	GuestName    string    `json:"guest_name" binding:"required"`
	GuestEmail   string    `json:"guest_email" binding:"required,email"`
	GuestPhone   string    `json:"guest_phone"`
	PartySize    int       `json:"party_size" binding:"required,min=1,max=20"`
	ReservedAt   time.Time `json:"reserved_at" binding:"required"`
	DurationMins int       `json:"duration_mins"`
	Notes        string    `json:"notes"`
}

type ReservationStatusUpdate struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
}

type TimeSlot struct {
	Time            string  `json:"time"`
	AvailableTables []Table `json:"available_tables"`
}

type AvailabilityResponse struct {
	Date           string     `json:"date"`
	AvailableSlots []TimeSlot `json:"available_slots"`
}
