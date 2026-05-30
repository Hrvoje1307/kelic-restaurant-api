package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
)

type AuthHandler struct {
	supabaseURL     string
	supabaseAnonKey string
}

func NewAuthHandler(supabaseURL, supabaseAnonKey string) *AuthHandler {
	return &AuthHandler{supabaseURL: supabaseURL, supabaseAnonKey: supabaseAnonKey}
}

// AdminLogin godoc
// @Summary      Admin login
// @Description  Authenticates an admin or superadmin and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.AdminLoginRequest   true  "Credentials"
// @Success      200    {object}  models.AdminLoginResponse
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var input models.AdminLoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, _ := json.Marshal(map[string]string{
		"email":    input.Email,
		"password": input.Password,
	})

	req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodPost,
		fmt.Sprintf("%s/auth/v1/token?grant_type=password", h.supabaseURL),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", h.supabaseAnonKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth service unavailable"})
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials", "code": "UNAUTHORIZED"})
		return
	}

	var supabaseResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		User        struct {
			UserMetadata struct {
				Role string `json:"role"`
			} `json:"user_metadata"`
		} `json:"user"`
	}
	if err := json.Unmarshal(raw, &supabaseResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse auth response"})
		return
	}

	role := supabaseResp.User.UserMetadata.Role
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "code": "FORBIDDEN"})
		return
	}

	c.JSON(http.StatusOK, models.AdminLoginResponse{
		AccessToken: supabaseResp.AccessToken,
		TokenType:   supabaseResp.TokenType,
		ExpiresIn:   supabaseResp.ExpiresIn,
	})
}
