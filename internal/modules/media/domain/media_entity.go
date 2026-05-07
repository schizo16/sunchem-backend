package domain

import "time"

type MediaFile struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FileName  string    `json:"name" gorm:"size:500"`
	FilePath  string    `json:"file_path" gorm:"size:1000"`
	FileSize  int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type IMediaRepository interface {
	FindAll() ([]MediaFile, error)
	Create(file *MediaFile) error
	Delete(id uint) error
}
