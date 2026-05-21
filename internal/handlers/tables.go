package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"
)

type TableHandler struct {
	repo *repository.TableRepo
}

func NewTableHandler(repo *repository.TableRepo) *TableHandler {
	return &TableHandler{repo: repo}
}

// List godoc
// @Summary      Dohvati sve stolove
// @Description  Javni endpoint
// @Tags         Stolovi
// @Produce      json
// @Success      200  {array}  models.Table
// @Router       /tables [get]
func (h *TableHandler) List(c *gin.Context) {
	list, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []models.Table{}
	}
	c.JSON(http.StatusOK, list)
}

// Create godoc
// @Summary      Dodaj stol (admin)
// @Tags         Stolovi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      models.TableInput  true  "Stol"
// @Success      201    {object}  models.Table
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Router       /tables [post]
func (h *TableHandler) Create(c *gin.Context) {
	var input models.TableInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// Update godoc
// @Summary      Ažuriraj stol (admin)
// @Tags         Stolovi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string             true  "UUID stola"
// @Param        input  body      models.TableInput  true  "Stol"
// @Success      200    {object}  models.Table
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /tables/{id} [put]
func (h *TableHandler) Update(c *gin.Context) {
	var input models.TableInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.repo.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stol nije pronađen", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// Delete godoc
// @Summary      Obriši stol (admin)
// @Tags         Stolovi
// @Security     BearerAuth
// @Param        id   path  string  true  "UUID stola"
// @Success      204
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /tables/{id} [delete]
func (h *TableHandler) Delete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stol nije pronađen", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
