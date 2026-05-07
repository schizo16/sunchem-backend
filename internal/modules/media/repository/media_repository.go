package repository

import (
	"sunchem-backend/internal/modules/media/domain"

	"gorm.io/gorm"
)

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) domain.IMediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) FindAll() ([]domain.MediaFile, error) {
	var files []domain.MediaFile
	err := r.db.Order("created_at desc").Find(&files).Error
	return files, err
}

func (r *mediaRepository) Create(file *domain.MediaFile) error {
	return r.db.Create(file).Error
}

func (r *mediaRepository) Delete(id uint) error {
	return r.db.Delete(&domain.MediaFile{}, id).Error
}
