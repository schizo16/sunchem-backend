package usecase

import (
	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/blog/domain"
)

type BlogUseCase struct {
	repo domain.IBlogRepository
}

func NewBlogUseCase(repo domain.IBlogRepository) *BlogUseCase {
	return &BlogUseCase{repo: repo}
}

func (uc *BlogUseCase) List() ([]domain.BlogPost, *errors.AppError) {
	posts, err := uc.repo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn bài viết")
	}
	return posts, nil
}

func (uc *BlogUseCase) GetByID(id uint) (*domain.BlogPost, *errors.AppError) {
	post, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, 404, "NOT_FOUND", "Bài viết không tồn tại")
	}
	post.Views++
	_ = uc.repo.Update(post)
	return post, nil
}

func (uc *BlogUseCase) GetBySlug(slug string) (*domain.BlogPost, *errors.AppError) {
	post, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return nil, errors.Wrap(err, 404, "NOT_FOUND", "Bài viết không tồn tại")
	}
	post.Views++
	_ = uc.repo.Update(post)
	return post, nil
}

func (uc *BlogUseCase) Create(post *domain.BlogPost) *errors.AppError {
	if post.Title == "" {
		return errors.ErrBadRequest
	}
	if err := uc.repo.Create(post); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi tạo bài viết")
	}
	return nil
}

func (uc *BlogUseCase) Update(id uint, post *domain.BlogPost) *errors.AppError {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.Wrap(err, 404, "NOT_FOUND", "Bài viết không tồn tại")
	}
	post.ID = id
	if err := uc.repo.Update(post); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi cập nhật bài viết")
	}
	return nil
}

func (uc *BlogUseCase) UpdatePartial(id uint, updates map[string]interface{}) (*domain.BlogPost, *errors.AppError) {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, 404, "NOT_FOUND", "Bài viết không tồn tại")
	}
	if err := uc.repo.UpdateFields(id, updates); err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi cập nhật bài viết")
	}
	post, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn sau cập nhật")
	}
	return post, nil
}

func (uc *BlogUseCase) Delete(id uint) *errors.AppError {
	if err := uc.repo.Delete(id); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi xóa bài viết")
	}
	return nil
}
