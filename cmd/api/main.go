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
	menuItemHandler := handlers.NewMenuItemHandler(repository.NewMenuItemRepo(pool))
	reservationHandler := handlers.NewReservationHandler(repository.NewReservationRepo(pool))
	tableHandler := handlers.NewTableHandler(repository.NewTableRepo(pool))

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	admin := middleware.RequireRole("admin", "superadmin")
	auth := middleware.RequireAuth()

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

			items := menu.Group("/items")
			{
				items.GET("", menuItemHandler.List)
				items.GET("/:id", menuItemHandler.GetByID)
				items.POST("", admin, menuItemHandler.Create)
				items.PUT("/:id", admin, menuItemHandler.Update)
				items.DELETE("/:id", admin, menuItemHandler.Delete)
			}
		}

		tables := api.Group("/tables")
		{
			tables.GET("", tableHandler.List)
			tables.POST("", admin, tableHandler.Create)
			tables.PUT("/:id", admin, tableHandler.Update)
			tables.DELETE("/:id", admin, tableHandler.Delete)
		}

		reservations := api.Group("/reservations")
		{
			reservations.GET("", admin, reservationHandler.List)
			reservations.POST("", reservationHandler.Create)
			reservations.GET("/availability", reservationHandler.Availability)
			reservations.GET("/my", auth, reservationHandler.MyReservations)
			reservations.GET("/:id", admin, reservationHandler.GetByID)
			reservations.PATCH("/:id", admin, reservationHandler.UpdateStatus)
			reservations.DELETE("/:id", auth, reservationHandler.Delete)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
