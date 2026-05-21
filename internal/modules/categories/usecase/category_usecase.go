package usecase

import (
	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/categories/domain"
)

type CategoryUseCase struct {
	repo domain.ICategoryRepository
}

func NewCategoryUseCase(repo domain.ICategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}

func (uc *CategoryUseCase) List() ([]domain.Category, *errors.AppError) {
	cats, err := uc.repo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn danh mục")
	}
	return cats, nil
}

func (uc *CategoryUseCase) Create(cat *domain.Category) *errors.AppError {
	if cat.Name == "" {
		return errors.ErrBadRequest
	}
	if err := uc.repo.Create(cat); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi tạo danh mục")
	}
	return nil
}

func (uc *CategoryUseCase) Update(id uint, cat *domain.Category) *errors.AppError {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.Wrap(err, 404, "NOT_FOUND", "Danh mục không tồn tại")
	}
	cat.ID = id
	if err := uc.repo.Update(cat); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi cập nhật danh mục")
	}
	return nil
}

func (uc *CategoryUseCase) Delete(id uint) *errors.AppError {
	if err := uc.repo.Delete(id); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi xóa danh mục")
	}
	return nil
}
