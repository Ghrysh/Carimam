package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ghrysh/carimam/order-service/internal/models"
	"github.com/ghrysh/carimam/order-service/internal/repository"
)

type OrderItemRequest struct {
	ProductID uint `json:"product_id" form:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" form:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" form:"items" binding:"required,min=1"`
}

type ProductResponse struct {
	Status string `json:"status"`
	Data   struct {
		Price float64 `json:"price"`
		Stock int     `json:"stock"`
	} `json:"data"`
}

type OrderUseCase interface {
	CreateOrder(eaterID uint, req CreateOrderRequest) error
	GetMyOrders(eaterID uint) ([]models.Order, error)
}

type orderUseCase struct {
	repo repository.OrderRepository
}

func NewOrderUseCase(repo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{repo}
}

func (u *orderUseCase) CreateOrder(eaterID uint, req CreateOrderRequest) error {
	var totalPrice float64
	var orderItems []models.OrderItem

	// =======================================================
	// 1. TELEPON PRODUCT SERVICE (Untuk cek harga & stok)
	// =======================================================
	for _, item := range req.Items {
		url := fmt.Sprintf("http://localhost:8081/api/products/%d", item.ProductID)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			return fmt.Errorf("gagal menghubungi layanan produk untuk makanan ID %d", item.ProductID)
		}
		
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var productData ProductResponse
		json.Unmarshal(body, &productData)

		if productData.Status != "success" {
			return fmt.Errorf("makanan ID %d tidak ditemukan", item.ProductID)
		}

		if productData.Data.Stock < item.Quantity {
			return fmt.Errorf("stok makanan ID %d tidak mencukupi", item.ProductID)
		}

		price := productData.Data.Price
		subTotal := price * float64(item.Quantity)
		totalPrice += subTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})
	}

	// =======================================================
	// 2. TELEPON IDENTITY SERVICE (Untuk potong saldo)
	// =======================================================
	deductPayload := map[string]interface{}{
		"user_id": eaterID,
		"amount":  totalPrice,
	}
	jsonData, _ := json.Marshal(deductPayload)

	identityURL := "http://localhost:8080/api/internal/users/deduct"
	respIdentity, err := http.Post(identityURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("gagal menghubungi layanan kasir pusat")
	}
	defer respIdentity.Body.Close()

	if respIdentity.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		json.NewDecoder(respIdentity.Body).Decode(&errorResponse)
		
		errorMessage := "Pembayaran gagal"
		if msg, ok := errorResponse["message"].(string); ok {
			errorMessage = msg
		}
		return errors.New(errorMessage)
	}

	// =======================================================
	// 3. SIMPAN KE DATABASE ORDER
	// =======================================================
	newOrder := &models.Order{
		EaterID:    eaterID,
		TotalPrice: totalPrice,
		Status:     "paid",
		Items:      orderItems,
	}

	return u.repo.CreateOrder(newOrder)
}

func (u *orderUseCase) GetMyOrders(eaterID uint) ([]models.Order, error) {
	return u.repo.GetByEaterID(eaterID)
}