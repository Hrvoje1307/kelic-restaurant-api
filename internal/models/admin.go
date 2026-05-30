package models

import "time"

type Admin struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// adminWithHash is used internally for login — password_hash is never serialised.
type AdminWithHash struct {
	Admin
	PasswordHash string `json:"-"`
}
