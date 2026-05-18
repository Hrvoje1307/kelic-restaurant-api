package main

import (
	"log"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/config"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/handlers"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		api.GET("/health", handlers.Health)

		menu := api.Group("/menu")
		{
			menu.GET("", handlers.ListMenuItems)
			menu.GET("/:id", handlers.GetMenuItem)
			menu.POST("", handlers.CreateMenuItem)
			menu.PUT("/:id", handlers.UpdateMenuItem)
			menu.DELETE("/:id", handlers.DeleteMenuItem)
		}

		orders := api.Group("/orders")
		{
			orders.GET("", handlers.ListOrders)
			orders.GET("/:id", handlers.GetOrder)
			orders.POST("", handlers.CreateOrder)
			orders.PUT("/:id/status", handlers.UpdateOrderStatus)
		}

		tables := api.Group("/tables")
		{
			tables.GET("", handlers.ListTables)
			tables.GET("/:id", handlers.GetTable)
			tables.POST("", handlers.CreateTable)
			tables.PUT("/:id", handlers.UpdateTable)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
