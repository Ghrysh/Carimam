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
		"message": "Pesanan berhasil dibuat! Menunggu pembayaran. 📝",
	})
}