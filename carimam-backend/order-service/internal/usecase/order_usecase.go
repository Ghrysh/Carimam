package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ghrysh/carimam/order-service/internal/models"
	"github.com/ghrysh/carimam/order-service/internal/repository"
)

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
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

	newOrder := &models.Order{
		EaterID:    eaterID,
		TotalPrice: totalPrice,
		Status:     "pending",
		Items:      orderItems,
	}

	return u.repo.CreateOrder(newOrder)
}