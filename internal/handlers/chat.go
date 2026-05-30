package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/cache"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"
)

type ChatHandler struct {
	repo         *repository.ChatRepo
	ai           *anthropic.Client
	systemPrompt string
	cache        *cache.ResponseCache
}

func NewChatHandler(repo *repository.ChatRepo, ai *anthropic.Client, systemPrompt string, c *cache.ResponseCache) *ChatHandler {
	return &ChatHandler{repo: repo, ai: ai, systemPrompt: systemPrompt, cache: c}
}

// Chat godoc
// @Summary      Pošalji poruku AI asistentu
// @Description  Streaming SSE endpoint — vraća Server-Sent Events
// @Tags         AI Chat
// @Accept       json
// @Produce      text/event-stream
// @Param        input  body  models.ChatRequest  true  "Poruka"
// @Success      200    {object}  models.ChatResponse
// @Failure      400    {object}  map[string]string
// @Router       /ai/chat [post]
func (h *ChatHandler) Chat(c *gin.Context) {
	var input models.ChatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	history, err := h.repo.GetHistory(c.Request.Context(), input.SessionID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if history == nil {
		history = []models.ChatMessage{}
	}

	history = append(history, models.ChatMessage{Role: "user", Content: input.Message})

	// For first-message questions (FAQs) return a cached answer without hitting the API.
	isFAQ := len(history) == 1
	cacheKey := cache.NormalizeKey(input.Message)
	if isFAQ {
		if cached, ok := h.cache.Get(cacheKey); ok {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("X-Accel-Buffering", "no")

			resp := models.ChatResponse{SessionID: input.SessionID, Message: cached, Done: false}
			data, _ := json.Marshal(resp)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)

			done := models.ChatResponse{SessionID: input.SessionID, Done: true}
			data, _ = json.Marshal(done)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()

			newHistory := append(history, models.ChatMessage{Role: "assistant", Content: cached})
			_ = h.repo.SaveMessages(context.Background(), input.SessionID, newHistory)
			return
		}
	}

	var msgParams []anthropic.MessageParam
	for _, m := range history {
		role := anthropic.MessageParamRoleUser
		if m.Role == "assistant" {
			role = anthropic.MessageParamRoleAssistant
		}
		msgParams = append(msgParams, anthropic.MessageParam{
			Role: role,
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{Text: m.Content},
			}},
		})
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	stream := h.ai.Messages.NewStreaming(c.Request.Context(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 1024,
		System:    []anthropic.TextBlockParam{{Text: h.systemPrompt}},
		Messages:  msgParams,
	})

	var fullResponse strings.Builder

	for stream.Next() {
		event := stream.Current()
		switch e := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch delta := e.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				fullResponse.WriteString(delta.Text)
				resp := models.ChatResponse{
					SessionID: input.SessionID,
					Message:   delta.Text,
					Done:      false,
				}
				data, _ := json.Marshal(resp)
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				c.Writer.Flush()
			}
		}
	}

	if err := stream.Err(); err != nil {
		errData, _ := json.Marshal(gin.H{"error": err.Error()})
		fmt.Fprintf(c.Writer, "data: %s\n\n", errData)
		c.Writer.Flush()
		return
	}

	doneResp := models.ChatResponse{
		SessionID: input.SessionID,
		Message:   "",
		Done:      true,
	}
	data, _ := json.Marshal(doneResp)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	c.Writer.Flush()

	answer := fullResponse.String()
	if isFAQ {
		h.cache.Set(cacheKey, answer)
	}
	history = append(history, models.ChatMessage{Role: "assistant", Content: answer})
	_ = h.repo.SaveMessages(context.Background(), input.SessionID, history)
}

// GetHistory godoc
// @Summary      Dohvati povijest razgovora
// @Tags         AI Chat
// @Produce      json
// @Param        session_id  path      string  true  "ID sesije"
// @Success      200         {object}  models.ChatSession
// @Failure      404         {object}  map[string]string
// @Router       /ai/chat/{session_id} [get]
func (h *ChatHandler) GetHistory(c *gin.Context) {
	sessionID := c.Param("session_id")
	messages, err := h.repo.GetHistory(c.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "sesija nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if messages == nil {
		messages = []models.ChatMessage{}
	}
	c.JSON(http.StatusOK, models.ChatSession{
		SessionID: sessionID,
		Messages:  messages,
	})
}

// DeleteSession godoc
// @Summary      Obriši sesiju chata
// @Tags         AI Chat
// @Param        session_id  path  string  true  "ID sesije"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /ai/chat/{session_id} [delete]
func (h *ChatHandler) DeleteSession(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("session_id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "sesija nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
