package main

import (
	"log"

	"github.com/hrvojecuckovic/kelic-restaurant/internal/config"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/handlers"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/middleware"
	_ "github.com/hrvojecuckovic/kelic-restaurant/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Kelic Restaurant API
// @version         1.0
// @description     Backend API for Kelic Restaurant
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	cfg := config.Load()

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.GET("/health", handlers.Health)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
