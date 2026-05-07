package domain

import "time"

type BlogPost struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"size:500"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;size:500"`
	Summary     string    `json:"summary" gorm:"type:text"`
	Content     string    `json:"content" gorm:"type:text"`
	Thumbnail   string    `json:"thumbnail" gorm:"size:1000"`
	Category    string    `json:"category" gorm:"size:200"`
	Status      string    `json:"status" gorm:"size:20;default:draft"`
	Views       int       `json:"views" gorm:"default:0"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IBlogRepository interface {
	FindAll() ([]BlogPost, error)
	FindByID(id uint) (*BlogPost, error)
	FindBySlug(slug string) (*BlogPost, error)
	Create(post *BlogPost) error
	Update(post *BlogPost) error
	UpdateFields(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}
