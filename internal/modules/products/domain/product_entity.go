package domain

import "time"

type Product struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Slug             string    `json:"slug" gorm:"uniqueIndex;size:300"`
	Name             string    `json:"name" gorm:"size:500"`
	ShortDescription string    `json:"short_description" gorm:"type:text"`
	Image            string    `json:"image" gorm:"size:1000"`
	Category         string    `json:"category" gorm:"size:200"`
	Highlights       string    `json:"highlights" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type IProductRepository interface {
	FindAll() ([]Product, error)
	FindByID(id uint) (*Product, error)
	FindBySlug(slug string) (*Product, error)
	Create(product *Product) error
	Update(product *Product) error
	UpdateFields(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}
