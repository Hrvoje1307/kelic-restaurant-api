package models

import "time"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type ChatResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	Done      bool   `json:"done"`
}

type ChatSession struct {
	SessionID string        `json:"session_id"`
	Messages  []ChatMessage `json:"messages"`
}

type ChatSessionFull struct {
	SessionID string        `json:"session_id"`
	Messages  []ChatMessage `json:"messages"`
	UpdatedAt time.Time     `json:"updated_at"`
}
