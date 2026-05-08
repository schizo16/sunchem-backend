package domain

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	GenoractID *string   `json:"genoract_id" gorm:"uniqueIndex;size:255"`
	Username   string    `json:"username" gorm:"uniqueIndex;size:100"`
	Password   string    `json:"-" gorm:"size:255"`
	Name       string    `json:"name" gorm:"size:200"`
	Role       string    `json:"role" gorm:"size:20;default:employee"`
	CreatedAt  time.Time `json:"created_at"`
}

type IUserRepository interface {
	FindByUsername(username string) (*User, error)
	FindByGenoractID(genoractID string) (*User, error)
	FindByID(id uint) (*User, error)
	FindAll() ([]User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}
