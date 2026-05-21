package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"
)

type UserHandler struct {
	repo *repository.UserRepo
}

func NewUserHandler(repo *repository.UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetMe godoc
// @Summary      Dohvati vlastiti profil
// @Tags         Korisnici
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserProfile
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, ok := userID.(string)
	if !ok || id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not identify user"})
		return
	}
	u, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "profil nije pronađen", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

// UpdateMe godoc
// @Summary      Ažuriraj vlastiti profil
// @Tags         Korisnici
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.UserProfileInput  true  "Profil"
// @Success      200    {object}  models.UserProfile
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Router       /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, ok := userID.(string)
	if !ok || id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not identify user"})
		return
	}
	var input models.UserProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.repo.Upsert(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

// List godoc
// @Summary      Dohvati sve korisnike (admin)
// @Tags         Korisnici
// @Produce      json
// @Security     BearerAuth
// @Param        role  query  string  false  "Filter po ulozi"  Enums(guest,admin,superadmin)
// @Success      200   {array}   models.UserProfile
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	list, err := h.repo.List(c.Request.Context(), c.Query("role"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []models.UserProfile{}
	}
	c.JSON(http.StatusOK, list)
}

// UpdateRole godoc
// @Summary      Promijeni ulogu korisnika (superadmin)
// @Tags         Korisnici
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                true  "UUID korisnika"
// @Param        input  body      models.UserRoleUpdate  true  "Nova uloga"
// @Success      200    {object}  models.UserProfile
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /users/{id}/role [patch]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	var input models.UserRoleUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.repo.UpdateRole(c.Request.Context(), c.Param("id"), input.Role)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "korisnik nije pronađen", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}
