package usecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/media/domain"
)

type MediaUseCase struct {
	repo      domain.IMediaRepository
	uploadDir string
}

func NewMediaUseCase(repo domain.IMediaRepository, uploadDir string) *MediaUseCase {
	os.MkdirAll(uploadDir, 0755)
	return &MediaUseCase{repo: repo, uploadDir: uploadDir}
}

func (uc *MediaUseCase) List() ([]domain.MediaFile, *errors.AppError) {
	files, err := uc.repo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn media")
	}
	return files, nil
}

func (uc *MediaUseCase) Upload(fileName string, reader io.Reader, size int64) (*domain.MediaFile, *errors.AppError) {
	ext := filepath.Ext(fileName)
	storeName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	storePath := filepath.Join(uc.uploadDir, storeName)

	dst, err := os.Create(storePath)
	if err != nil {
		return nil, errors.Wrap(err, 500, "UPLOAD_ERROR", "Lỗi lưu file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return nil, errors.Wrap(err, 500, "UPLOAD_ERROR", "Lỗi ghi file")
	}

	file := &domain.MediaFile{
		FileName:  fileName,
		FilePath:  "/uploads/" + storeName,
		FileSize:  size,
		CreatedAt: time.Now(),
	}
	if err := uc.repo.Create(file); err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi lưu thông tin file")
	}
	return file, nil
}

func (uc *MediaUseCase) Delete(id uint) *errors.AppError {
	if err := uc.repo.Delete(id); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi xóa file")
	}
	return nil
}
