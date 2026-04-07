package repository

import (
	"github.com/ghrysh/carimam/order-service/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(eaterID uint) (*models.Cart, error)
	AddToCartItem(cartID uint, item *models.CartItem) error
	RemoveCartItem(cartID uint, productID uint) error
	ClearCart(cartID uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db}
}

func (r *cartRepository) GetCart(eaterID uint) (*models.Cart, error) {
	var cart models.Cart
	
	err := r.db.Preload("Items").Where("eater_id = ?", eaterID).First(&cart).Error

	if err != nil {
		cart = models.Cart{EaterID: eaterID}
		errCreate := r.db.Create(&cart).Error
		if errCreate != nil {
			return nil, errCreate
		}
	}
	
	return &cart, nil
}

func (r *cartRepository) AddToCartItem(cartID uint, item *models.CartItem) error {
	var existing models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, item.ProductID).First(&existing).Error
	if err == nil {
		existing.Quantity += item.Quantity
		return r.db.Save(&existing).Error
	}
	item.CartID = cartID
	return r.db.Create(item).Error
}

func (r *cartRepository) RemoveCartItem(cartID uint, productID uint) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}
