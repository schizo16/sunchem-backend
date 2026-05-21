package repository

import (
	"sunchem-backend/internal/modules/tags/domain"

	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) domain.ITagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll() ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.Order("created_at desc").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindByID(id uint) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *tagRepository) Create(tag *domain.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) Update(tag *domain.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Tag{}, id).Error
}
