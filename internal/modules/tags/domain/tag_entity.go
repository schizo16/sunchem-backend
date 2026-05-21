package domain

import "time"

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:300"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;size:300"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ITagRepository interface {
	FindAll() ([]Tag, error)
	FindByID(id uint) (*Tag, error)
	Create(tag *Tag) error
	Update(tag *Tag) error
	Delete(id uint) error
}
