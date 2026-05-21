package repository

import (
	"context"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

const userCols = `id::text, COALESCE(full_name,''), COALESCE(phone,''), COALESCE(role,'guest'), COALESCE(avatar_url,''), created_at, updated_at`

func scanUser(row interface{ Scan(...any) error }) (models.UserProfile, error) {
	var u models.UserProfile
	err := row.Scan(&u.ID, &u.FullName, &u.Phone, &u.Role, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (models.UserProfile, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+userCols+` FROM user_profiles WHERE id = $1`, id,
	)
	u, err := scanUser(row)
	return u, mapNotFound(err)
}

func (r *UserRepo) List(ctx context.Context, role string) ([]models.UserProfile, error) {
	query := `SELECT ` + userCols + ` FROM user_profiles`
	args := []any{}
	if role != "" {
		query += ` WHERE role = $1`
		args = append(args, role)
	}
	query += ` ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.UserProfile
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

func (r *UserRepo) Upsert(ctx context.Context, id string, input models.UserProfileInput) (models.UserProfile, error) {
	row := r.db.QueryRow(ctx,
		`INSERT INTO user_profiles (id, full_name, phone, avatar_url)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (id) DO UPDATE
		   SET full_name  = EXCLUDED.full_name,
		       phone      = EXCLUDED.phone,
		       avatar_url = EXCLUDED.avatar_url,
		       updated_at = NOW()
		 RETURNING `+userCols,
		id, input.FullName, input.Phone, input.AvatarURL,
	)
	return scanUser(row)
}

func (r *UserRepo) UpdateRole(ctx context.Context, id string, role string) (models.UserProfile, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE user_profiles SET role = $1, updated_at = NOW()
		 WHERE id = $2
		 RETURNING `+userCols,
		role, id,
	)
	u, err := scanUser(row)
	return u, mapNotFound(err)
}
