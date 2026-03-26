package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ghrysh/carimam/identity-service/internal/config"
	"github.com/ghrysh/carimam/identity-service/internal/models"

	"github.com/ghrysh/carimam/identity-service/internal/repository"
	"github.com/ghrysh/carimam/identity-service/internal/usecase"
	handler "github.com/ghrysh/carimam/identity-service/internal/delivery/http"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Peringatan: file .env tidak ditemukan")
	}

	db := config.SetupDatabase()

	db.AutoMigrate(&models.User{})

	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	
	userUseCase := usecase.NewUserUseCase(userRepo)
	
	handler.NewUserHandler(r, userUseCase)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong! CariMam is Running 🚀"})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}