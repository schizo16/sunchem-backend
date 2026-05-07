package repository

import (
	"time"

	"sunchem-backend/internal/modules/analytics/domain"

	"gorm.io/gorm"
)

type trafficRepository struct {
	db *gorm.DB
}

func NewTrafficRepository(db *gorm.DB) domain.ITrafficRepository {
	return &trafficRepository{db: db}
}

func (r *trafficRepository) Create(event *domain.TrafficEvent) error {
	return r.db.Create(event).Error
}

func (r *trafficRepository) CountSince(since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&domain.TrafficEvent{}).Where("timestamp >= ?", since).Count(&count).Error
	return count, err
}

func (r *trafficRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&domain.TrafficEvent{}).Count(&count).Error
	return count, err
}

func (r *trafficRepository) FindAll(offset, limit int) ([]domain.TrafficEvent, error) {
	var events []domain.TrafficEvent
	err := r.db.Order("timestamp desc").Offset(offset).Limit(limit).Find(&events).Error
	return events, err
}

func (r *trafficRepository) AvgDuration() (float64, error) {
	var avg float64
	err := r.db.Model(&domain.TrafficEvent{}).Select("COALESCE(AVG(duration), 0)").Scan(&avg).Error
	return avg, err
}

func (r *trafficRepository) BounceRate() (float64, error) {
	var total int64
	r.db.Model(&domain.TrafficEvent{}).Count(&total)
	if total == 0 {
		return 0, nil
	}
	var bounceCount int64
	r.db.Model(&domain.TrafficEvent{}).Where("duration <= 10").Count(&bounceCount)
	return float64(bounceCount) / float64(total) * 100, nil
}

func (r *trafficRepository) PageViewCounts() (map[string]int64, error) {
	type row struct {
		Page  string
		Count int64
	}
	var rows []row
	err := r.db.Model(&domain.TrafficEvent{}).Select("page, COUNT(*) as count").Group("page").Order("count desc").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]int64)
	for _, r := range rows {
		result[r.Page] = r.Count
	}
	return result, nil
}
