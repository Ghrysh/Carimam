package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductReview struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	EaterID   uint           `gorm:"not null;index" json:"eater_id"`
	Rating    int            `gorm:"not null" json:"rating"`
	Comment   string         `gorm:"type:text" json:"comment"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}