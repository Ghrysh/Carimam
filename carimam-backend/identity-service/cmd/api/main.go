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
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"status": "error", "message": "Sesi tidak valid"})
			return
		}

		var id uint
		switch v := userID.(type) {
		case float64:
			id = uint(v)
		case uint:
			id = v
		default:
			c.JSON(500, gin.H{"status": "error", "message": "Tipe data ID tidak dikenali"})
			return
		}

		user, err := userUseCase.GetProfile(id)
		if err != nil {
			c.JSON(404, gin.H{"status": "error", "message": "User tidak ditemukan"})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Berhasil mengambil data profil",
			"data": gin.H{
				"user_id": user.ID,
				"name":    user.Name,
				"email":   user.Email,
				"role":    user.Role,
				"balance": user.Balance,
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