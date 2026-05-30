package repository

import (
	"context"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) GetByEmail(ctx context.Context, email string) (models.AdminWithHash, error) {
	var a models.AdminWithHash
	err := r.db.QueryRow(ctx,
		`SELECT id::text, email, password_hash, COALESCE(full_name,''), role, created_at, updated_at
		 FROM admins WHERE email = $1`, email,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.FullName, &a.Role, &a.CreatedAt, &a.UpdatedAt)
	return a, mapNotFound(err)
}
