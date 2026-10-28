package main

import (
	"log"

	"backend/config"
	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/repository"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg.DBFile)

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(cfg.CORSOrigin))

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", authHandler.Profile)
		}
	}

	log.Printf("Server listening on port :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
