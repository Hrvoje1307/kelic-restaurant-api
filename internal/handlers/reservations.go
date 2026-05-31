package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"
)

type ReservationHandler struct {
	repo *repository.ReservationRepo
}

func NewReservationHandler(repo *repository.ReservationRepo) *ReservationHandler {
	return &ReservationHandler{repo: repo}
}

func (h *ReservationHandler) List(c *gin.Context) {
	filter := repository.ReservationFilter{
		Date:    c.Query("date"),
		Status:  c.Query("status"),
		TableID: c.Query("table_id"),
	}
	list, err := h.repo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []models.Reservation{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var input models.ReservationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.repo.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *ReservationHandler) Availability(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}
	partySize := 0
	if _, err := fmt.Sscanf(c.Query("party_size"), "%d", &partySize); err != nil || partySize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "party_size must be a positive integer"})
		return
	}

	slots, err := repository.TimeSlots(date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var availableSlots []models.TimeSlot
	for _, slot := range slots {
		tables, err := h.repo.AvailableTablesForSlot(c.Request.Context(), slot, 90, partySize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if tables == nil {
			tables = []models.Table{}
		}
		availableSlots = append(availableSlots, models.TimeSlot{
			Time:            slot.Format("15:04"),
			AvailableTables: tables,
		})
	}

	c.JSON(http.StatusOK, models.AvailabilityResponse{
		Date:           date,
		AvailableSlots: availableSlots,
	})
}

func (h *ReservationHandler) MyReservations(c *gin.Context) {
	email, _ := c.Get("email")
	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not identify user"})
		return
	}
	list, err := h.repo.GetByGuestEmail(c.Request.Context(), emailStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []models.Reservation{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *ReservationHandler) GetByID(c *gin.Context) {
	res, err := h.repo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "rezervacija nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ReservationHandler) UpdateStatus(c *gin.Context) {
	var input models.ReservationStatusUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.repo.UpdateStatus(c.Request.Context(), c.Param("id"), input.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "rezervacija nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ReservationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	if role != "admin" && role != "superadmin" {
		// guest — verify ownership
		email, _ := c.Get("email")
		res, err := h.repo.GetByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "rezervacija nije pronađena", "code": "NOT_FOUND"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if res.GuestEmail != email {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "code": "FORBIDDEN"})
			return
		}
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "rezervacija nije pronađena", "code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
