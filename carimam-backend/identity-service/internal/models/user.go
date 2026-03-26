package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Phone     string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`

	Role      string         `gorm:"type:varchar(20);default:'eater'" json:"role"` 

	Address   string         `gorm:"type:text" json:"address"`
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	
	Balance   float64        `gorm:"default:0" json:"balance"`
	Rating    float32        `gorm:"default:0" json:"rating"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}