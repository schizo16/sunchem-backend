package usecase

import (
	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/tags/domain"
)

type TagUseCase struct {
	repo domain.ITagRepository
}

func NewTagUseCase(repo domain.ITagRepository) *TagUseCase {
	return &TagUseCase{repo: repo}
}

func (uc *TagUseCase) List() ([]domain.Tag, *errors.AppError) {
	tags, err := uc.repo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn thẻ")
	}
	return tags, nil
}

func (uc *TagUseCase) Create(tag *domain.Tag) *errors.AppError {
	if tag.Name == "" {
		return errors.ErrBadRequest
	}
	if err := uc.repo.Create(tag); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi tạo thẻ")
	}
	return nil
}

func (uc *TagUseCase) Update(id uint, tag *domain.Tag) *errors.AppError {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.Wrap(err, 404, "NOT_FOUND", "Thẻ không tồn tại")
	}
	tag.ID = id
	if err := uc.repo.Update(tag); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi cập nhật thẻ")
	}
	return nil
}

func (uc *TagUseCase) Delete(id uint) *errors.AppError {
	if err := uc.repo.Delete(id); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi xóa thẻ")
	}
	return nil
}
