package repository

import (
	"github.com/ghrysh/carimam/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetByEaterID(eaterID uint) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByEaterID(eaterID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items").Where("eater_id = ?", eaterID).Order("created_at desc").Find(&orders).Error
	return orders, err
}