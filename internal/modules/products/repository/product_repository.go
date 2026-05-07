package repository

import (
	"sunchem-backend/internal/modules/products/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.IProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Order("id asc").Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id uint) (*domain.Product, error) {
	var p domain.Product
	err := r.db.First(&p, id).Error
	return &p, err
}

func (r *productRepository) FindBySlug(slug string) (*domain.Product, error) {
	var p domain.Product
	err := r.db.Where("slug = ?", slug).First(&p).Error
	return &p, err
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).Updates(updates).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}
