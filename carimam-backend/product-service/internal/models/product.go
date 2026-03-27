package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	
	CookID      uint           `gorm:"not null;index" json:"cook_id"` 
	
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Category    string         `gorm:"type:varchar(100)" json:"category"`
	ImageURL    string         `gorm:"type:varchar(255)" json:"image_url"`
	
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
	Stock       int            `gorm:"default:0" json:"stock"`

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}