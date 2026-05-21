package domain

import "time"

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:300"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;size:300"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ICategoryRepository interface {
	FindAll() ([]Category, error)
	FindByID(id uint) (*Category, error)
	Create(cat *Category) error
	Update(cat *Category) error
	Delete(id uint) error
}
