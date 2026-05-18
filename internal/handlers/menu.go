package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
)

// In-memory store — replace with DB repository
var menuItems = []models.MenuItem{
	{ID: 1, Name: "Bruschetta", Description: "Toasted bread with tomatoes", Price: 7.50, Category: models.CategoryStarter, Available: true},
	{ID: 2, Name: "Margherita Pizza", Description: "Classic tomato and mozzarella", Price: 12.00, Category: models.CategoryMain, Available: true},
	{ID: 3, Name: "Tiramisu", Description: "Classic Italian dessert", Price: 6.00, Category: models.CategoryDessert, Available: true},
}
var menuNextID = 4

func ListMenuItems(c *gin.Context) {
	c.JSON(http.StatusOK, menuItems)
}

func GetMenuItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	for _, item := range menuItems {
		if item.ID == id {
			c.JSON(http.StatusOK, item)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
}

func CreateMenuItem(c *gin.Context) {
	var input models.MenuItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = menuNextID
	menuNextID++
	menuItems = append(menuItems, input)
	c.JSON(http.StatusCreated, input)
}

func UpdateMenuItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input models.MenuItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, item := range menuItems {
		if item.ID == id {
			input.ID = id
			menuItems[i] = input
			c.JSON(http.StatusOK, input)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
}

func DeleteMenuItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	for i, item := range menuItems {
		if item.ID == id {
			menuItems = append(menuItems[:i], menuItems[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
}
