package repository

import (
	"sunchem-backend/internal/modules/blog/domain"

	"gorm.io/gorm"
)

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) domain.IBlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) FindAll() ([]domain.BlogPost, error) {
	var posts []domain.BlogPost
	err := r.db.Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *blogRepository) FindByID(id uint) (*domain.BlogPost, error) {
	var post domain.BlogPost
	err := r.db.First(&post, id).Error
	return &post, err
}

func (r *blogRepository) FindBySlug(slug string) (*domain.BlogPost, error) {
	var post domain.BlogPost
	err := r.db.Where("slug = ?", slug).First(&post).Error
	return &post, err
}

func (r *blogRepository) Create(post *domain.BlogPost) error {
	return r.db.Create(post).Error
}

func (r *blogRepository) Update(post *domain.BlogPost) error {
	return r.db.Save(post).Error
}

func (r *blogRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&domain.BlogPost{}).Where("id = ?", id).Updates(updates).Error
}

func (r *blogRepository) Delete(id uint) error {
	return r.db.Delete(&domain.BlogPost{}, id).Error
}
