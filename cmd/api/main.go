package main

import (
	"log"

	_ "github.com/hrvojecuckovic/kelic-restaurant/docs"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/config"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/db"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/handlers"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/middleware"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Kelic Restaurant API
// @version         1.0
// @description     Backend API for Kelic Restaurant
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	pool, err := db.Connect(cfg.DSN)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	menuCategoryHandler := handlers.NewMenuCategoryHandler(repository.NewMenuCategoryRepo(pool))

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	admin := middleware.RequireRole("admin", "superadmin")

	api := r.Group("/api/v1")
	{
		api.GET("/health", handlers.Health)

		menu := api.Group("/menu")
		{
			categories := menu.Group("/categories")
			{
				categories.GET("", menuCategoryHandler.List)
				categories.POST("", admin, menuCategoryHandler.Create)
				categories.PUT("/:id", admin, menuCategoryHandler.Update)
				categories.DELETE("/:id", admin, menuCategoryHandler.Delete)
			}
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
