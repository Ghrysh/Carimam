package http

import (
	"net/http"

	"github.com/ghrysh/carimam/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase usecase.OrderUseCase
}

func NewOrderHandler(r *gin.Engine, usecase usecase.OrderUseCase, eaterMiddleware gin.HandlerFunc) {
	handler := &OrderHandler{usecase}

	eaterGroup := r.Group("/api")
	eaterGroup.Use(eaterMiddleware)
	{
		eaterGroup.POST("/orders", handler.CreateOrder)
		eaterGroup.GET("/orders", handler.GetMyOrders)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req usecase.CreateOrderRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format keranjang belanja salah"})
		return
	}

	eaterID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Sesi tidak valid"})
		return
	}

	err := h.usecase.CreateOrder(eaterID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibuat dan saldo otomatis terpotong! 💸",
	})
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	eaterID, _ := c.Get("user_id")

	orders, err := h.usecase.GetMyOrders(eaterID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil riwayat pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil riwayat pesanan",
		"data":    orders,
	})
}