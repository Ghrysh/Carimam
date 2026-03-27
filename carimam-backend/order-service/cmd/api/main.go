package main

import (
	"log"
	"os"

	"github.com/ghrysh/carimam/order-service/internal/config"
	"github.com/ghrysh/carimam/order-service/internal/models"
	"github.com/ghrysh/carimam/order-service/internal/middleware"
	"github.com/ghrysh/carimam/order-service/internal/repository"
	"github.com/ghrysh/carimam/order-service/internal/usecase"
	handler "github.com/ghrysh/carimam/order-service/internal/delivery/http"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Peringatan: file .env tidak ditemukan")
	}

	db := config.SetupDatabase()

	log.Println("Menjalankan proses migrasi database order...")
	errMigrate := db.AutoMigrate(&models.Order{}, &models.OrderItem{})
	if errMigrate != nil {
		log.Fatalf("Gagal melakukan migrasi: %v", errMigrate)
	}
	log.Println("Migrasi tabel Order dan OrderItem selesai! ✅")

	r := gin.Default()

	eaterOnlyMiddleware := middleware.AuthMiddleware("eater")

	orderRepo := repository.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUseCase(orderRepo)
	handler.NewOrderHandler(r, orderUseCase, eaterOnlyMiddleware)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong! CariMam Order Service is Running 🛒🚀"})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}