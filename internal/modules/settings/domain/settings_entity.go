package domain

import "time"

type Setting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex;size:100"`
	Value     string    `json:"value" gorm:"type:text"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ISettingRepository interface {
	FindByKey(key string) (*Setting, error)
	FindAll() ([]Setting, error)
	Upsert(setting *Setting) error
}
