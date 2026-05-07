package usecase

import (
	"time"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/settings/domain"
)

type SettingUseCase struct {
	repo domain.ISettingRepository
}

func NewSettingUseCase(repo domain.ISettingRepository) *SettingUseCase {
	return &SettingUseCase{repo: repo}
}

func (uc *SettingUseCase) GetAll() (map[string]string, *errors.AppError) {
	settings, err := uc.repo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn")
	}
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

func (uc *SettingUseCase) Save(values map[string]string) *errors.AppError {
	for k, v := range values {
		s := &domain.Setting{
			Key:       k,
			Value:     v,
			UpdatedAt: time.Now(),
		}
		if err := uc.repo.Upsert(s); err != nil {
			return errors.Wrap(err, 500, "DB_ERROR", "Lỗi lưu cài đặt")
		}
	}
	return nil
}
