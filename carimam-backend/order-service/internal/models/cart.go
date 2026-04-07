package models

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	EaterID   uint       `gorm:"uniqueIndex;not null" json:"eater_id"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items"`
	
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	CartID    uint `gorm:"not null;index" json:"cart_id"`
	ProductID uint `gorm:"not null" json:"product_id"`
	Quantity  int  `gorm:"not null" json:"quantity"`
}