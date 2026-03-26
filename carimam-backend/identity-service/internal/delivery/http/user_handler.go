package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ghrysh/carimam/identity-service/internal/usecase"
)

type UserHandler struct {
	usecase usecase.UserUseCase
}

func NewUserHandler(r *gin.Engine, usecase usecase.UserUseCase) {
	handler := &UserHandler{usecase}

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak lengkap atau format salah",
			"error":   err.Error(),
		})
		return
	}

	err := h.usecase.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Yeay! Akun berhasil didaftarkan. Silakan login.",
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak lengkap atau format salah",
			"error":   err.Error(),
		})
		return
	}

	token, err := h.usecase.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil!",
		"data": gin.H{
			"token": token,
		},
	})
}