package repository

import (
	"sunchem-backend/internal/modules/settings/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) domain.ISettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) FindByKey(key string) (*domain.Setting, error) {
	var s domain.Setting
	err := r.db.Where("key = ?", key).First(&s).Error
	return &s, err
}

func (r *settingRepository) FindAll() ([]domain.Setting, error) {
	var settings []domain.Setting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *settingRepository) Upsert(setting *domain.Setting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(setting).Error
}
