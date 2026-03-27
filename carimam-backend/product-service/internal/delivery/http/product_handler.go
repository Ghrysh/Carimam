package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ghrysh/carimam/product-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase usecase.ProductUseCase
}

func NewProductHandler(r *gin.Engine, usecase usecase.ProductUseCase, cookMiddleware gin.HandlerFunc) {
	handler := &ProductHandler{usecase}

	r.GET("/api/products", handler.GetAllProducts)
	r.GET("/api/products/:id", handler.GetProductByID)

	cookGroup := r.Group("/api")
	cookGroup.Use(cookMiddleware)
	{
		cookGroup.POST("/products", handler.CreateProduct)
		cookGroup.PUT("/products/:id", handler.UpdateProduct)
		cookGroup.DELETE("/products/:id", handler.DeleteProduct)
		cookGroup.PATCH("/products/:id/image", handler.UploadImage)
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req usecase.CreateProductRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data makanan tidak lengkap atau format salah",
			"error":   err.Error(),
		})
		return
	}

	cookID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Sesi tidak valid"})
		return
	}

	productID, err := h.usecase.CreateProduct(cookID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan menu makanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Menu makanan baru berhasil ditambahkan! 🍳",
		"data": gin.H{
			"id": productID,
		},
	})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.usecase.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar makanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil daftar makanan",
		"data":    products,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))
	cookID, _ := c.Get("user_id")

	var req usecase.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format data salah",
			"error":   err.Error(),
		})
		return
	}

	err := h.usecase.UpdateProduct(cookID.(uint), uint(productID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Menu berhasil diupdate! 📝",
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))
	cookID, _ := c.Get("user_id")

	err := h.usecase.DeleteProduct(cookID.(uint), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Menu berhasil dihapus! 🗑️",
	})
}

func (h *ProductHandler) UploadImage(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))
	cookID, _ := c.Get("user_id")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "File gambar tidak ditemukan"})
		return
	}

	uploadDir := "uploads/products"
	os.MkdirAll(uploadDir, os.ModePerm)

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("carimam_product_%d_%d%s", productID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan gambar ke server"})
		return
	}

	imageURL := fmt.Sprintf("http://localhost:8081/%s", filePath)

	err = h.usecase.UpdateProductImage(cookID.(uint), uint(productID), imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Foto makanan berhasil diunggah! 📸",
		"data": gin.H{
			"image_url": imageURL,
		},
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))

	product, err := h.usecase.GetProductByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Menu makanan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    product,
	})
}