package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ghrysh/carimam/product-service/internal/config"
	handler "github.com/ghrysh/carimam/product-service/internal/delivery/http"
	"github.com/ghrysh/carimam/product-service/internal/middleware"
	"github.com/ghrysh/carimam/product-service/internal/models"
	"github.com/ghrysh/carimam/product-service/internal/repository"
	"github.com/ghrysh/carimam/product-service/internal/usecase"
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

	db.AutoMigrate(&models.Product{}, &models.ProductReview{})

	// ==========================================================
	// INIT GIN ENGINE & WIRING (Merakit Komponen)
	// ==========================================================
	r := gin.Default()

	cookOnlyMiddleware := middleware.AuthMiddleware("cook")
	eaterMiddleware := middleware.AuthMiddleware("eater")

	productRepo := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepo)

	handler.NewProductHandler(r, productUseCase, cookOnlyMiddleware, eaterMiddleware)

	// ==========================================================
	// PUBLIC ROUTES
	// ==========================================================
	r.Static("/uploads", "./uploads")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong! CariMam Product Service is Running 🍔🚀"})
	})

	// ==========================================================
	// START SERVER
	// ==========================================================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}