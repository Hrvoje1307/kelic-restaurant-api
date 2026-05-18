package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
)

var orders []models.Order
var orderNextID = 1

func ListOrders(c *gin.Context) {
	c.JSON(http.StatusOK, orders)
}

func GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	for _, o := range orders {
		if o.ID == id {
			c.JSON(http.StatusOK, o)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}

func CreateOrder(c *gin.Context) {
	var input struct {
		TableID int               `json:"table_id" binding:"required"`
		Items   []models.OrderItem `json:"items" binding:"required,min=1"`
		Notes   string            `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total := 0.0
	for _, item := range input.Items {
		total += item.UnitPrice * float64(item.Quantity)
	}

	order := models.Order{
		ID:         orderNextID,
		TableID:    input.TableID,
		Items:      input.Items,
		Status:     models.OrderStatusPending,
		TotalPrice: total,
		Notes:      input.Notes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	orderNextID++
	orders = append(orders, order)
	c.JSON(http.StatusCreated, order)
}

func UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, o := range orders {
		if o.ID == id {
			orders[i].Status = input.Status
			orders[i].UpdatedAt = time.Now()
			c.JSON(http.StatusOK, orders[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
}
