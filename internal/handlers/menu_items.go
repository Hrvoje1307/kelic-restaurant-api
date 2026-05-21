package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"
)

type MenuItemHandler struct {
	repo *repository.MenuItemRepo
}

func NewMenuItemHandler(repo *repository.MenuItemRepo) *MenuItemHandler {
	return &MenuItemHandler{repo: repo}
}

// List godoc
// @Summary      Dohvati sve stavke menija
// @Description  Javni endpoint. Filtriranje po kategoriji i dostupnosti.
// @Tags         Meni
// @Produce      json
// @Param        category_id     query  string  false  "UUID kategorije"
// @Param        available_only  query  bool    false  "Samo dostupne stavke"  default(true)
// @Success      200  {array}   models.MenuItem
// @Router       /menu/items [get]
func (h *MenuItemHandler) List(c *gin.Context) {
	filter := repository.MenuItemFilter{
		AvailableOnly: c.DefaultQuery("available_only", "true") == "true",
	}
	if catID := c.Query("category_id"); catID != "" {
		filter.CategoryID = &catID
	}

	items, err := h.repo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []models.MenuItem{}
	}
	c.JSON(http.StatusOK, items)
}

// GetByID godoc
// @Summary      Dohvati jednu stavku
// @Tags         Meni
// @Produce      json
// @Param        id   path      string  true  "UUID stavke"
// @Success      200  {object}  models.MenuItem
// @Failure      404  {object}  map[string]string
// @Router       /menu/items/{id} [get]
func (h *MenuItemHandler) GetByID(c *gin.Context) {
	item, err := h.repo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stavka nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Create godoc
// @Summary      Dodaj stavku menija (admin)
// @Tags         Meni
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.MenuItemInput  true  "Stavka menija"
// @Success      201    {object}  models.MenuItem
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Router       /menu/items [post]
func (h *MenuItemHandler) Create(c *gin.Context) {
	var input models.MenuItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// Update godoc
// @Summary      Ažuriraj stavku (admin)
// @Tags         Meni
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                true  "UUID stavke"
// @Param        input  body      models.MenuItemInput  true  "Stavka menija"
// @Success      200    {object}  models.MenuItem
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /menu/items/{id} [put]
func (h *MenuItemHandler) Update(c *gin.Context) {
	var input models.MenuItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.repo.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stavka nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Delete godoc
// @Summary      Obriši stavku (admin)
// @Tags         Meni
// @Security     BearerAuth
// @Param        id   path  string  true  "UUID stavke"
// @Success      204
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /menu/items/{id} [delete]
func (h *MenuItemHandler) Delete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stavka nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
