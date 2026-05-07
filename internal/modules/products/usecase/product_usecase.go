package usecase

import (
	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/products/domain"
)

type ProductUseCase struct {
	repo domain.IProductRepository
}

func NewProductUseCase(repo domain.IProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (uc *ProductUseCase) List() ([]domain.Product, *errors.AppError) {
	products, err := uc.repo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn sản phẩm")
	}
	return products, nil
}

func (uc *ProductUseCase) GetByID(id uint) (*domain.Product, *errors.AppError) {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, 404, "NOT_FOUND", "Sản phẩm không tồn tại")
	}
	return p, nil
}

func (uc *ProductUseCase) GetBySlug(slug string) (*domain.Product, *errors.AppError) {
	p, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return nil, errors.Wrap(err, 404, "NOT_FOUND", "Sản phẩm không tồn tại")
	}
	return p, nil
}

func (uc *ProductUseCase) Create(product *domain.Product) *errors.AppError {
	if product.Name == "" || product.Slug == "" {
		return errors.ErrBadRequest
	}
	if err := uc.repo.Create(product); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi tạo sản phẩm")
	}
	return nil
}

func (uc *ProductUseCase) UpdatePartial(id uint, updates map[string]interface{}) (*domain.Product, *errors.AppError) {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, 404, "NOT_FOUND", "Sản phẩm không tồn tại")
	}
	if err := uc.repo.UpdateFields(id, updates); err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi cập nhật sản phẩm")
	}
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn sau cập nhật")
	}
	return p, nil
}

func (uc *ProductUseCase) Delete(id uint) *errors.AppError {
	if err := uc.repo.Delete(id); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi xóa sản phẩm")
	}
	return nil
}
