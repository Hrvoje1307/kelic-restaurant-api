package repository

import (
	"context"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TableRepo struct {
	db *pgxpool.Pool
}

func NewTableRepo(db *pgxpool.Pool) *TableRepo {
	return &TableRepo{db: db}
}

const tableCols = `id, table_number, capacity, COALESCE(location, ''), is_active, created_at`

func scanTable(row interface{ Scan(...any) error }) (models.Table, error) {
	var t models.Table
	err := row.Scan(&t.ID, &t.TableNumber, &t.Capacity, &t.Location, &t.IsActive, &t.CreatedAt)
	return t, err
}

func (r *TableRepo) List(ctx context.Context) ([]models.Table, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+tableCols+` FROM tables ORDER BY table_number ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Table
	for rows.Next() {
		t, err := scanTable(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, nil
}

func (r *TableRepo) Create(ctx context.Context, input models.TableInput) (models.Table, error) {
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}
	row := r.db.QueryRow(ctx,
		`INSERT INTO tables (table_number, capacity, location, is_active)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+tableCols,
		input.TableNumber, input.Capacity, input.Location, isActive,
	)
	return scanTable(row)
}

func (r *TableRepo) Update(ctx context.Context, id string, input models.TableInput) (models.Table, error) {
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}
	row := r.db.QueryRow(ctx,
		`UPDATE tables SET table_number=$1, capacity=$2, location=$3, is_active=$4
		 WHERE id=$5
		 RETURNING `+tableCols,
		input.TableNumber, input.Capacity, input.Location, isActive, id,
	)
	t, err := scanTable(row)
	return t, mapNotFound(err)
}

func (r *TableRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM tables WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
