package usecase

import (
	"errors"

	"github.com/ghrysh/carimam/order-service/internal/models"
	"github.com/ghrysh/carimam/order-service/internal/repository"
)

type AddCartRequest struct {
	ProductID uint `json:"product_id" form:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" form:"quantity" binding:"required,gt=0"`
}

type CartUseCase interface {
	GetCart(eaterID uint) (*models.Cart, error)
	AddToCart(eaterID uint, req AddCartRequest) error
	RemoveFromCart(eaterID uint, productID uint) error
	CheckoutCart(eaterID uint) error
}

type cartUseCase struct {
	cartRepo     repository.CartRepository
	orderUseCase OrderUseCase
}

func NewCartUseCase(cartRepo repository.CartRepository, orderUseCase OrderUseCase) CartUseCase {
	return &cartUseCase{cartRepo, orderUseCase}
}

func (u *cartUseCase) GetCart(eaterID uint) (*models.Cart, error) {
	return u.cartRepo.GetCart(eaterID)
}

func (u *cartUseCase) AddToCart(eaterID uint, req AddCartRequest) error {
	cart, err := u.cartRepo.GetCart(eaterID)
	if err != nil {
		return err
	}
	item := &models.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	return u.cartRepo.AddToCartItem(cart.ID, item)
}

func (u *cartUseCase) RemoveFromCart(eaterID uint, productID uint) error {
	cart, err := u.cartRepo.GetCart(eaterID)
	if err != nil {
		return err
	}
	return u.cartRepo.RemoveCartItem(cart.ID, productID)
}

func (u *cartUseCase) CheckoutCart(eaterID uint) error {
	cart, err := u.cartRepo.GetCart(eaterID)
	if err != nil || len(cart.Items) == 0 {
		return errors.New("keranjang belanja masih kosong nih")
	}

	var orderItems []OrderItemRequest
	for _, item := range cart.Items {
		orderItems = append(orderItems, OrderItemRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	err = u.orderUseCase.CreateOrder(eaterID, CreateOrderRequest{Items: orderItems})
	if err != nil {
		return err
	}

	return u.cartRepo.ClearCart(cart.ID)
}
