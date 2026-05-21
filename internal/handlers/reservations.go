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

// List godoc
// @Summary      Dohvati rezervacije (admin)
// @Tags         Rezervacije
// @Produce      json
// @Security     BearerAuth
// @Param        date      query  string  false  "Filter po datumu (YYYY-MM-DD)"
// @Param        status    query  string  false  "Filter po statusu"  Enums(pending,confirmed,cancelled,completed)
// @Param        table_id  query  string  false  "Filter po stolu (UUID)"
// @Success      200  {array}   models.Reservation
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /reservations [get]
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

// Create godoc
// @Summary      Kreiraj rezervaciju
// @Description  Javni endpoint — gost ne treba biti prijavljen
// @Tags         Rezervacije
// @Accept       json
// @Produce      json
// @Param        input  body      models.ReservationInput  true  "Rezervacija"
// @Success      201    {object}  models.Reservation
// @Failure      400    {object}  map[string]string
// @Failure      409    {object}  map[string]string
// @Router       /reservations [post]
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

// Availability godoc
// @Summary      Provjeri slobodne termine
// @Description  Javni endpoint — za prikaz dostupnih termina u formi
// @Tags         Rezervacije
// @Produce      json
// @Param        date        query  string  true  "Datum (YYYY-MM-DD)"
// @Param        party_size  query  int     true  "Broj gostiju"
// @Success      200  {object}  models.AvailabilityResponse
// @Failure      400  {object}  map[string]string
// @Router       /reservations/availability [get]
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

// MyReservations godoc
// @Summary      Moje rezervacije (gost)
// @Tags         Rezervacije
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Reservation
// @Failure      401  {object}  map[string]string
// @Router       /reservations/my [get]
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

// GetByID godoc
// @Summary      Dohvati rezervaciju (admin)
// @Tags         Rezervacije
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "UUID rezervacije"
// @Success      200  {object}  models.Reservation
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /reservations/{id} [get]
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

// UpdateStatus godoc
// @Summary      Promijeni status rezervacije (admin)
// @Tags         Rezervacije
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                          true  "UUID rezervacije"
// @Param        input  body      models.ReservationStatusUpdate  true  "Novi status"
// @Success      200    {object}  models.Reservation
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /reservations/{id} [patch]
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

// Delete godoc
// @Summary      Otkaži rezervaciju
// @Description  Gost može otkazati svoju, admin može otkazati bilo koju
// @Tags         Rezervacije
// @Security     BearerAuth
// @Param        id   path  string  true  "UUID rezervacije"
// @Success      204
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /reservations/{id} [delete]
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
