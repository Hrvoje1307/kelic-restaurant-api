package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/models"
)

var tables = []models.Table{
	{ID: 1, Number: 1, Capacity: 2, Status: models.TableStatusFree},
	{ID: 2, Number: 2, Capacity: 4, Status: models.TableStatusFree},
	{ID: 3, Number: 3, Capacity: 6, Status: models.TableStatusFree},
}
var tableNextID = 4

func ListTables(c *gin.Context) {
	c.JSON(http.StatusOK, tables)
}

func GetTable(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	for _, t := range tables {
		if t.ID == id {
			c.JSON(http.StatusOK, t)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
}

func CreateTable(c *gin.Context) {
	var input models.Table
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = tableNextID
	tableNextID++
	if input.Status == "" {
		input.Status = models.TableStatusFree
	}
	tables = append(tables, input)
	c.JSON(http.StatusCreated, input)
}

func UpdateTable(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input models.Table
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, t := range tables {
		if t.ID == id {
			input.ID = id
			tables[i] = input
			c.JSON(http.StatusOK, input)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
}
