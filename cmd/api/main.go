package main

import (
	"log"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	_ "github.com/hrvojecuckovic/kelic-restaurant/docs"
	"github.com/hrvojecuckovic/kelic-restaurant/internal/cache"
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
// @host            kelic-restaurant-api-production.up.railway.app
// @BasePath        /api/v1
// @schemes         https
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

	aiClient := anthropic.NewClient(option.WithAPIKey(cfg.AnthropicAPIKey))

	promptBytes, err := os.ReadFile("prompts/starter.txt")
	if err != nil {
		log.Fatalf("Failed to read system prompt: %v", err)
	}

	menuCategoryHandler := handlers.NewMenuCategoryHandler(repository.NewMenuCategoryRepo(pool))
	menuItemHandler := handlers.NewMenuItemHandler(repository.NewMenuItemRepo(pool))
	reservationHandler := handlers.NewReservationHandler(repository.NewReservationRepo(pool))
	tableHandler := handlers.NewTableHandler(repository.NewTableRepo(pool))
	responseCache := cache.New(24 * time.Hour)
	chatHandler := handlers.NewChatHandler(repository.NewChatRepo(pool), &aiClient, string(promptBytes), responseCache)
	userHandler := handlers.NewUserHandler(repository.NewUserRepo(pool))

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	admin := middleware.RequireRole("admin", "superadmin")
	superadmin := middleware.RequireRole("superadmin")
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

		ai := api.Group("/ai")
		{
			ai.GET("/chats", admin, chatHandler.ListSessions)
			chat := ai.Group("/chat")
			{
				chat.POST("", chatHandler.Chat)
				chat.GET("/:session_id", chatHandler.GetHistory)
				chat.DELETE("/:session_id", chatHandler.DeleteSession)
			}
		}

		users := api.Group("/users")
		{
			users.GET("/me", auth, userHandler.GetMe)
			users.PUT("/me", auth, userHandler.UpdateMe)
			users.GET("", admin, userHandler.List)
			users.PATCH("/:id/role", superadmin, userHandler.UpdateRole)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
