package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"
)

type MenuCategoryHandler struct {
	repo *repository.MenuCategoryRepo
}

func NewMenuCategoryHandler(repo *repository.MenuCategoryRepo) *MenuCategoryHandler {
	return &MenuCategoryHandler{repo: repo}
}

// List godoc
// @Summary      Dohvati sve kategorije
// @Description  Javni endpoint — ne treba autentikaciju
// @Tags         Meni
// @Produce      json
// @Success      200  {array}   models.MenuCategory
// @Router       /menu/categories [get]
func (h *MenuCategoryHandler) List(c *gin.Context) {
	cats, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cats == nil {
		cats = []models.MenuCategory{}
	}
	c.JSON(http.StatusOK, cats)
}

// Create godoc
// @Summary      Kreiraj kategoriju (admin)
// @Tags         Meni
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.MenuCategoryInput  true  "Kategorija"
// @Success      201    {object}  models.MenuCategory
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Router       /menu/categories [post]
func (h *MenuCategoryHandler) Create(c *gin.Context) {
	var input models.MenuCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// Update godoc
// @Summary      Ažuriraj kategoriju (admin)
// @Tags         Meni
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                    true  "UUID kategorije"
// @Param        input  body      models.MenuCategoryInput  true  "Kategorija"
// @Success      200    {object}  models.MenuCategory
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /menu/categories/{id} [put]
func (h *MenuCategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input models.MenuCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat, err := h.repo.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kategorija nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// Delete godoc
// @Summary      Obriši kategoriju (admin)
// @Tags         Meni
// @Security     BearerAuth
// @Param        id   path  string  true  "UUID kategorije"
// @Success      204
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /menu/categories/{id} [delete]
func (h *MenuCategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kategorija nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
