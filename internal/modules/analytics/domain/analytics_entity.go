package domain

import "time"

type TrafficEvent struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SessionID string    `json:"session_id" gorm:"index;size:36"`
	Page      string    `json:"page" gorm:"index;size:500"`
	IP        string    `json:"ip" gorm:"size:45"`
	Location  string    `json:"location" gorm:"size:255"`
	Device    string    `json:"device" gorm:"size:100"`
	Duration  int       `json:"duration"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

type TrafficSummary struct {
	Today       int64 `json:"today"`
	Week        int64 `json:"week"`
	Month       int64 `json:"month"`
	PageViews   int64 `json:"pageViews"`
	AvgDuration int64 `json:"avgDuration"`
	BounceRate  float64 `json:"bounceRate"`
}

type ITrafficRepository interface {
	Create(event *TrafficEvent) error
	CountSince(since time.Time) (int64, error)
	CountAll() (int64, error)
	FindAll(offset, limit int) ([]TrafficEvent, error)
	AvgDuration() (float64, error)
	BounceRate() (float64, error)
	PageViewCounts() (map[string]int64, error)
}
