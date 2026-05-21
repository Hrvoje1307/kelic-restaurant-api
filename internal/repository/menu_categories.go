package repository

import (
	"context"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuCategoryRepo struct {
	db *pgxpool.Pool
}

func NewMenuCategoryRepo(db *pgxpool.Pool) *MenuCategoryRepo {
	return &MenuCategoryRepo{db: db}
}

func (r *MenuCategoryRepo) List(ctx context.Context) ([]models.MenuCategory, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, description, sort_order, created_at FROM menu_categories ORDER BY sort_order ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.MenuCategory
	for rows.Next() {
		var c models.MenuCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *MenuCategoryRepo) Create(ctx context.Context, input models.MenuCategoryInput) (models.MenuCategory, error) {
	var c models.MenuCategory
	err := r.db.QueryRow(ctx,
		`INSERT INTO menu_categories (name, description, sort_order)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, description, sort_order, created_at`,
		input.Name, input.Description, input.SortOrder,
	).Scan(&c.ID, &c.Name, &c.Description, &c.SortOrder, &c.CreatedAt)
	return c, err
}

func (r *MenuCategoryRepo) Update(ctx context.Context, id string, input models.MenuCategoryInput) (models.MenuCategory, error) {
	var c models.MenuCategory
	err := r.db.QueryRow(ctx,
		`UPDATE menu_categories
		 SET name = $1, description = $2, sort_order = $3
		 WHERE id = $4
		 RETURNING id, name, description, sort_order, created_at`,
		input.Name, input.Description, input.SortOrder, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.SortOrder, &c.CreatedAt)
	if err != nil {
		return c, mapNotFound(err)
	}
	return c, nil
}

func (r *MenuCategoryRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM menu_categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
