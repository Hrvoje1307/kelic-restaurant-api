package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepo struct {
	db *pgxpool.Pool
}

func NewReservationRepo(db *pgxpool.Pool) *ReservationRepo {
	return &ReservationRepo{db: db}
}

type ReservationFilter struct {
	Date    string
	Status  string
	TableID string
}

const reservationCols = `id, table_id, guest_name, guest_email, guest_phone,
	party_size, reserved_at, duration_mins, status, notes, created_at, updated_at`

func scanReservation(row interface{ Scan(...any) error }) (models.Reservation, error) {
	var r models.Reservation
	err := row.Scan(
		&r.ID, &r.TableID, &r.GuestName, &r.GuestEmail, &r.GuestPhone,
		&r.PartySize, &r.ReservedAt, &r.DurationMins, &r.Status, &r.Notes,
		&r.CreatedAt, &r.UpdatedAt,
	)
	return r, err
}

func (r *ReservationRepo) List(ctx context.Context, f ReservationFilter) ([]models.Reservation, error) {
	query := `SELECT ` + reservationCols + ` FROM reservations WHERE 1=1`
	args := []any{}
	i := 1

	if f.Date != "" {
		query += fmt.Sprintf(" AND reserved_at::date = $%d", i)
		args = append(args, f.Date)
		i++
	}
	if f.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, f.Status)
		i++
	}
	if f.TableID != "" {
		query += fmt.Sprintf(" AND table_id = $%d", i)
		args = append(args, f.TableID)
		i++
	}
	query += " ORDER BY reserved_at ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Reservation
	for rows.Next() {
		res, err := scanReservation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, res)
	}
	return list, nil
}

func (r *ReservationRepo) GetByID(ctx context.Context, id string) (models.Reservation, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+reservationCols+` FROM reservations WHERE id = $1`, id,
	)
	res, err := scanReservation(row)
	return res, mapNotFound(err)
}

func (r *ReservationRepo) GetByGuestEmail(ctx context.Context, email string) ([]models.Reservation, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+reservationCols+` FROM reservations WHERE guest_email = $1 ORDER BY reserved_at DESC`,
		email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Reservation
	for rows.Next() {
		res, err := scanReservation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, res)
	}
	return list, nil
}

func (r *ReservationRepo) Create(ctx context.Context, input models.ReservationInput) (models.Reservation, error) {
	if input.DurationMins == 0 {
		input.DurationMins = 90
	}
	row := r.db.QueryRow(ctx,
		`INSERT INTO reservations (table_id, guest_name, guest_email, guest_phone, party_size, reserved_at, duration_mins, notes)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING `+reservationCols,
		input.TableID, input.GuestName, input.GuestEmail, input.GuestPhone,
		input.PartySize, input.ReservedAt, input.DurationMins, input.Notes,
	)
	return scanReservation(row)
}

func (r *ReservationRepo) UpdateStatus(ctx context.Context, id, status string) (models.Reservation, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE reservations SET status = $1 WHERE id = $2 RETURNING `+reservationCols,
		status, id,
	)
	res, err := scanReservation(row)
	return res, mapNotFound(err)
}

func (r *ReservationRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM reservations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// AvailableTablesForSlot returns tables with enough capacity that have no overlapping reservation.
func (r *ReservationRepo) AvailableTablesForSlot(ctx context.Context, slotStart time.Time, durationMins, minCapacity int) ([]models.Table, error) {
	slotEnd := slotStart.Add(time.Duration(durationMins) * time.Minute)

	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.table_number, t.capacity, COALESCE(t.location, ''), t.is_active
		FROM tables t
		WHERE t.is_active = true
		  AND t.capacity >= $1
		  AND t.id NOT IN (
		    SELECT DISTINCT r.table_id
		    FROM reservations r
		    WHERE r.table_id IS NOT NULL
		      AND r.status NOT IN ('cancelled', 'completed')
		      AND r.reserved_at < $3
		      AND (r.reserved_at + (r.duration_mins * interval '1 minute')) > $2
		  )
		ORDER BY t.capacity ASC
	`, minCapacity, slotStart, slotEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var t models.Table
		if err := rows.Scan(&t.ID, &t.TableNumber, &t.Capacity, &t.Location, &t.IsActive); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

// TimeSlots returns hourly slots for a restaurant day (12:00–21:00).
func TimeSlots(date string) ([]time.Time, error) {
	loc, _ := time.LoadLocation("Europe/Zagreb")
	d, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	var slots []time.Time
	for h := 12; h <= 21; h++ {
		slots = append(slots, d.Add(time.Duration(h)*time.Hour))
	}
	return slots, nil
}
