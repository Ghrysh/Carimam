package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ghrysh/carimam/identity-service/internal/config"
	handler "github.com/ghrysh/carimam/identity-service/internal/delivery/http"
	"github.com/ghrysh/carimam/identity-service/internal/middleware"
	"github.com/ghrysh/carimam/identity-service/internal/models"
	"github.com/ghrysh/carimam/identity-service/internal/repository"
	"github.com/ghrysh/carimam/identity-service/internal/usecase"
)

func main() {
	// ==========================================================
	// INIT CONFIG & DATABASE
	// ==========================================================
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Peringatan: file .env tidak ditemukan")
	}

	db := config.SetupDatabase()

	db.AutoMigrate(&models.User{})

	// ==========================================================
	// INIT GIN ENGINE & WIRING
	// ==========================================================
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	handler.NewUserHandler(r, userUseCase)

	// ==========================================================
	// PUBLIC ROUTES
	// ==========================================================
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong! CariMam is Running 🚀"})
	})

	// ==========================================================
	// PROTECTED ROUTES (WITH MIDDLEWARE)
	// ==========================================================
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/profile", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Selamat datang di area rahasia!",
			"data": gin.H{
				"user_id": userID,
				"role":    role,
			},
		})
	})

	// ==========================================================
	// START SERVER
	// ==========================================================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}