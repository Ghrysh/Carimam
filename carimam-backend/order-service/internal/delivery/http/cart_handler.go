package http

import (
	"net/http"
	"strconv"

	"github.com/ghrysh/carimam/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	usecase usecase.CartUseCase
}

func NewCartHandler(r *gin.Engine, usecase usecase.CartUseCase, eaterMiddleware gin.HandlerFunc) {
	handler := &CartHandler{usecase}

	cartGroup := r.Group("/api/cart")
	cartGroup.Use(eaterMiddleware)
	{
		cartGroup.GET("", handler.GetCart)
		cartGroup.POST("", handler.AddToCart)
		cartGroup.DELETE("/:product_id", handler.RemoveFromCart)
		cartGroup.POST("/checkout", handler.CheckoutCart)
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	eaterID, _ := c.Get("user_id")
	cart, err := h.usecase.GetCart(eaterID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil keranjang"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": cart})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	var req usecase.AddCartRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format data salah"})
		return
	}
	eaterID, _ := c.Get("user_id")
	err := h.usecase.AddToCart(eaterID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Berhasil ditambahkan ke keranjang 🛒"})
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("product_id"))
	eaterID, _ := c.Get("user_id")

	err := h.usecase.RemoveFromCart(eaterID.(uint), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Berhasil dihapus dari keranjang"})
}

func (h *CartHandler) CheckoutCart(c *gin.Context) {
	eaterID, _ := c.Get("user_id")

	err := h.usecase.CheckoutCart(eaterID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Checkout dari keranjang berhasil! Saldo otomatis terpotong 💸"})
}