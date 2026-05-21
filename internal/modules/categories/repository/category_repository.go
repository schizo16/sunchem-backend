package repository

import (
	"sunchem-backend/internal/modules/categories/domain"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.ICategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]domain.Category, error) {
	var cats []domain.Category
	err := r.db.Order("created_at desc").Find(&cats).Error
	return cats, err
}

func (r *categoryRepository) FindByID(id uint) (*domain.Category, error) {
	var cat domain.Category
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *categoryRepository) Create(cat *domain.Category) error {
	return r.db.Create(cat).Error
}

func (r *categoryRepository) Update(cat *domain.Category) error {
	return r.db.Save(cat).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}
