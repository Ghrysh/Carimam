package repository

import (
	"github.com/ghrysh/carimam/product-service/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetAll() ([]models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	AddReview(review *models.ProductReview) error
	GetReviewsByProductID(productID uint) ([]models.ProductReview, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("is_available = ?", true).Find(&products).Error
	return products, err
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) AddReview(review *models.ProductReview) error {
	return r.db.Create(review).Error
}

func (r *productRepository) GetReviewsByProductID(productID uint) ([]models.ProductReview, error) {
	var reviews []models.ProductReview
	err := r.db.Where("product_id = ?", productID).Order("created_at desc").Find(&reviews).Error
	return reviews, err
}