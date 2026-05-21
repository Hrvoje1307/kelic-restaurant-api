package repository

import (
	"context"
	"fmt"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuItemRepo struct {
	db *pgxpool.Pool
}

func NewMenuItemRepo(db *pgxpool.Pool) *MenuItemRepo {
	return &MenuItemRepo{db: db}
}

type MenuItemFilter struct {
	CategoryID    *string
	AvailableOnly bool
}

func (r *MenuItemRepo) List(ctx context.Context, f MenuItemFilter) ([]models.MenuItem, error) {
	query := `SELECT id, category_id, name, description, price, image_url, is_available, allergens, tags, sort_order, created_at, updated_at
	          FROM menu_items WHERE 1=1`
	args := []any{}
	i := 1

	if f.AvailableOnly {
		query += fmt.Sprintf(" AND is_available = $%d", i)
		args = append(args, true)
		i++
	}
	if f.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", i)
		args = append(args, *f.CategoryID)
		i++
	}
	query += " ORDER BY sort_order ASC, name ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(
			&item.ID, &item.CategoryID, &item.Name, &item.Description,
			&item.Price, &item.ImageURL, &item.IsAvailable,
			&item.Allergens, &item.Tags, &item.SortOrder,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MenuItemRepo) GetByID(ctx context.Context, id string) (models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.QueryRow(ctx,
		`SELECT id, category_id, name, description, price, image_url, is_available, allergens, tags, sort_order, created_at, updated_at
		 FROM menu_items WHERE id = $1`, id,
	).Scan(
		&item.ID, &item.CategoryID, &item.Name, &item.Description,
		&item.Price, &item.ImageURL, &item.IsAvailable,
		&item.Allergens, &item.Tags, &item.SortOrder,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return item, mapNotFound(err)
	}
	return item, nil
}

func (r *MenuItemRepo) Create(ctx context.Context, input models.MenuItemInput) (models.MenuItem, error) {
	isAvailable := true
	if input.IsAvailable != nil {
		isAvailable = *input.IsAvailable
	}
	allergens := input.Allergens
	if allergens == nil {
		allergens = []string{}
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	var item models.MenuItem
	err := r.db.QueryRow(ctx,
		`INSERT INTO menu_items (category_id, name, description, price, image_url, is_available, allergens, tags, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, category_id, name, description, price, image_url, is_available, allergens, tags, sort_order, created_at, updated_at`,
		input.CategoryID, input.Name, input.Description, input.Price,
		input.ImageURL, isAvailable, allergens, tags, input.SortOrder,
	).Scan(
		&item.ID, &item.CategoryID, &item.Name, &item.Description,
		&item.Price, &item.ImageURL, &item.IsAvailable,
		&item.Allergens, &item.Tags, &item.SortOrder,
		&item.CreatedAt, &item.UpdatedAt,
	)
	return item, err
}

func (r *MenuItemRepo) Update(ctx context.Context, id string, input models.MenuItemInput) (models.MenuItem, error) {
	isAvailable := true
	if input.IsAvailable != nil {
		isAvailable = *input.IsAvailable
	}
	allergens := input.Allergens
	if allergens == nil {
		allergens = []string{}
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	var item models.MenuItem
	err := r.db.QueryRow(ctx,
		`UPDATE menu_items
		 SET category_id=$1, name=$2, description=$3, price=$4, image_url=$5,
		     is_available=$6, allergens=$7, tags=$8, sort_order=$9
		 WHERE id=$10
		 RETURNING id, category_id, name, description, price, image_url, is_available, allergens, tags, sort_order, created_at, updated_at`,
		input.CategoryID, input.Name, input.Description, input.Price,
		input.ImageURL, isAvailable, allergens, tags, input.SortOrder, id,
	).Scan(
		&item.ID, &item.CategoryID, &item.Name, &item.Description,
		&item.Price, &item.ImageURL, &item.IsAvailable,
		&item.Allergens, &item.Tags, &item.SortOrder,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return item, mapNotFound(err)
	}
	return item, nil
}

func (r *MenuItemRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM menu_items WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
