package repository

import (
	"context"
	"encoding/json"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepo struct {
	db *pgxpool.Pool
}

func NewChatRepo(db *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) GetHistory(ctx context.Context, sessionID string) ([]models.ChatMessage, error) {
	var raw []byte
	err := r.db.QueryRow(ctx,
		`SELECT messages FROM chat_sessions WHERE session_id = $1`,
		sessionID,
	).Scan(&raw)
	if err != nil {
		return nil, mapNotFound(err)
	}
	var messages []models.ChatMessage
	if err := json.Unmarshal(raw, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *ChatRepo) SaveMessages(ctx context.Context, sessionID string, messages []models.ChatMessage) error {
	raw, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx,
		`INSERT INTO chat_sessions (session_id, messages)
		 VALUES ($1, $2)
		 ON CONFLICT (session_id) DO UPDATE SET messages = $2, updated_at = NOW()`,
		sessionID, raw,
	)
	return err
}

func (r *ChatRepo) ListAll(ctx context.Context) ([]models.ChatSessionFull, error) {
	rows, err := r.db.Query(ctx,
		`SELECT session_id, messages, updated_at FROM chat_sessions ORDER BY updated_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.ChatSessionFull
	for rows.Next() {
		var s models.ChatSessionFull
		var raw []byte
		if err := rows.Scan(&s.SessionID, &raw, &s.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &s.Messages); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	if sessions == nil {
		sessions = []models.ChatSessionFull{}
	}
	return sessions, nil
}

func (r *ChatRepo) Delete(ctx context.Context, sessionID string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM chat_sessions WHERE session_id = $1`, sessionID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
