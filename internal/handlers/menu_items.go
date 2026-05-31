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
