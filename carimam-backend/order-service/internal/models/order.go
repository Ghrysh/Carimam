package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	EaterID    uint        `gorm:"not null;index" json:"eater_id"`
	TotalPrice float64     `gorm:"not null" json:"total_price"`
	Status     string      `gorm:"type:varchar(20);default:'pending'" json:"status"`

	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items"`

	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null;index" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"`
}