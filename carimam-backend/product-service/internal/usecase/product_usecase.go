package usecase

import (
	"errors"

	"github.com/ghrysh/carimam/product-service/internal/models"
	"github.com/ghrysh/carimam/product-service/internal/repository"
)

type CreateProductRequest struct {
	Name        string  `json:"name" form:"name" binding:"required"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
	Category    string  `json:"category" form:"category" binding:"required"`
	Stock       int     `json:"stock" form:"stock" binding:"gte=0"`
	ImageURL    string  `json:"image_url" form:"image_url"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

// ==========================================
// INTERFACE: DAFTAR KONTRAK FUNGSI
// ==========================================
type ProductUseCase interface {
	CreateProduct(cookID uint, req CreateProductRequest) (uint, error)
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	UpdateProduct(cookID uint, productID uint, req CreateProductRequest) error
	DeleteProduct(cookID uint, productID uint) error
	UpdateProductImage(cookID uint, productID uint, imageURL string) error
	AddProductReview(eaterID uint, productID uint, req CreateReviewRequest) error
	GetProductReviews(productID uint) ([]models.ProductReview, error)
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo}
}

// ==========================================
// IMPLEMENTASI FUNGSI
// ==========================================

func (u *productUseCase) CreateProduct(cookID uint, req CreateProductRequest) (uint, error) {
	newProduct := &models.Product{
		CookID:      cookID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsAvailable: true,
	}

	err := u.repo.Create(newProduct)
	return newProduct.ID, err 
}

func (u *productUseCase) GetAllProducts() ([]models.Product, error) {
	return u.repo.GetAll()
}

func (u *productUseCase) UpdateProduct(cookID uint, productID uint, req CreateProductRequest) error {
	product, err := u.repo.GetByID(productID)
	if err != nil {
		return errors.New("menu makanan tidak ditemukan")
	}

	if product.CookID != cookID {
		return errors.New("akses ditolak: kamu tidak bisa mengedit menu warung lain")
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Category = req.Category
	product.Stock = req.Stock
	product.ImageURL = req.ImageURL

	return u.repo.Update(product)
}

func (u *productUseCase) DeleteProduct(cookID uint, productID uint) error {
	product, err := u.repo.GetByID(productID)
	if err != nil {
		return errors.New("menu makanan tidak ditemukan")
	}

	if product.CookID != cookID {
		return errors.New("akses ditolak: kamu tidak bisa menghapus menu warung lain")
	}

	return u.repo.Delete(productID)
}

func (u *productUseCase) UpdateProductImage(cookID uint, productID uint, imageURL string) error {
	product, err := u.repo.GetByID(productID)
	if err != nil {
		return errors.New("menu makanan tidak ditemukan")
	}

	if product.CookID != cookID {
		return errors.New("akses ditolak: kamu tidak bisa mengganti foto menu warung lain")
	}

	product.ImageURL = imageURL
	return u.repo.Update(product)
}

func (u *productUseCase) GetProductByID(id uint) (*models.Product, error) {
	return u.repo.GetByID(id)
}

func (u *productUseCase) AddProductReview(eaterID uint, productID uint, req CreateReviewRequest) error {
	_, err := u.repo.GetByID(productID)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	review := &models.ProductReview{
		ProductID: productID,
		EaterID:   eaterID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	return u.repo.AddReview(review)
}

func (u *productUseCase) GetProductReviews(productID uint) ([]models.ProductReview, error) {
	return u.repo.GetReviewsByProductID(productID)
}