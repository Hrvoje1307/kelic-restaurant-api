package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo      *repository.AdminRepo
	jwtSecret []byte
}

func NewAuthHandler(repo *repository.AdminRepo, jwtSecret string) *AuthHandler {
	return &AuthHandler{repo: repo, jwtSecret: []byte(jwtSecret)}
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

	admin, err := h.repo.GetByEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials", "code": "UNAUTHORIZED"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials", "code": "UNAUTHORIZED"})
		return
	}

	expiresIn := 24 * time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   admin.ID,
		"email": admin.Email,
		"role":  admin.Role,
		"exp":   time.Now().Add(expiresIn).Unix(),
		"iat":   time.Now().Unix(),
	})

	signed, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, models.AdminLoginResponse{
		AccessToken: signed,
		TokenType:   "bearer",
		ExpiresIn:   int(expiresIn.Seconds()),
	})
}
