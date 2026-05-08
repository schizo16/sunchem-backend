package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	fmt.Println("[DB] Running auto migrations...")
	if err := db.AutoMigrate(
		&User{},
		&BlogPost{},
		&MediaFile{},
		&Setting{},
		&TrafficEvent{},
		&Product{},
	); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}
	fmt.Println("[DB] Migrations completed")
	return nil
}

type User struct {
	ID         uint      `gorm:"primaryKey"`
	GenoractID *string   `gorm:"uniqueIndex;size:255"`
	Username   string    `gorm:"uniqueIndex;size:100"`
	Password  string    `gorm:"size:255"`
	Name      string    `gorm:"size:200"`
	Role      string    `gorm:"size:20;default:employee"`
	CreatedAt time.Time
}

type BlogPost struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:500"`
	Slug        string    `gorm:"uniqueIndex;size:500"`
	Summary     string    `gorm:"type:text"`
	Content     string    `gorm:"type:text"`
	Thumbnail   string    `gorm:"size:1000"`
	Category    string    `gorm:"size:200"`
	Status      string    `gorm:"size:20;default:draft"`
	Views       int       `gorm:"default:0"`
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MediaFile struct {
	ID        uint      `gorm:"primaryKey"`
	FileName  string    `gorm:"size:500"`
	FilePath  string    `gorm:"size:1000"`
	FileSize  int64     `gorm:"default:0"`
	CreatedAt time.Time
}

type Setting struct {
	ID        uint      `gorm:"primaryKey"`
	Key       string    `gorm:"uniqueIndex;size:100"`
	Value     string    `gorm:"type:text"`
	UpdatedAt time.Time
}

type TrafficEvent struct {
	ID        uint      `gorm:"primaryKey"`
	SessionID string    `gorm:"index;size:36"`
	Page      string    `gorm:"index;size:500"`
	IP        string    `gorm:"size:45"`
	Location  string    `gorm:"size:255"`
	Device    string    `gorm:"size:100"`
	Duration  int
	Timestamp time.Time `gorm:"index"`
	CreatedAt time.Time
}

type Product struct {
	ID               uint      `gorm:"primaryKey"`
	Slug             string    `gorm:"uniqueIndex;size:300"`
	Name             string    `gorm:"size:500"`
	ShortDescription string    `gorm:"type:text"`
	Image            string    `gorm:"size:1000"`
	Category         string    `gorm:"size:200"`
	Highlights       string    `gorm:"type:text"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
