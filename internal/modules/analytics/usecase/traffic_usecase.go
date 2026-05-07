package usecase

import (
	"time"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/analytics/domain"
)

type TrafficUseCase struct {
	repo domain.ITrafficRepository
}

func NewTrafficUseCase(repo domain.ITrafficRepository) *TrafficUseCase {
	return &TrafficUseCase{repo: repo}
}

func (uc *TrafficUseCase) Track(event *domain.TrafficEvent) *errors.AppError {
	event.Timestamp = time.Now()
	if event.IP == "" {
		event.IP = "unknown"
	}
	if err := uc.repo.Create(event); err != nil {
		return errors.Wrap(err, 500, "DB_ERROR", "Lỗi ghi traffic")
	}
	return nil
}

func (uc *TrafficUseCase) Summary() (*domain.TrafficSummary, *errors.AppError) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	today, _ := uc.repo.CountSince(todayStart)
	week, _ := uc.repo.CountSince(weekStart)
	month, _ := uc.repo.CountSince(monthStart)
	total, _ := uc.repo.CountAll()
	avgDur, _ := uc.repo.AvgDuration()
	bounce, _ := uc.repo.BounceRate()

	return &domain.TrafficSummary{
		Today:       today,
		Week:        week,
		Month:       month,
		PageViews:   total,
		AvgDuration: int64(avgDur),
		BounceRate:  float64(int(bounce*100)) / 100,
	}, nil
}

func (uc *TrafficUseCase) Detail(page, limit int) ([]domain.TrafficEvent, *errors.AppError) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	events, err := uc.repo.FindAll(offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn traffic")
	}
	if events == nil {
		events = []domain.TrafficEvent{}
	}
	return events, nil
}

func (uc *TrafficUseCase) PageViews() (map[string]int64, *errors.AppError) {
	counts, err := uc.repo.PageViewCounts()
	if err != nil {
		return nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn")
	}
	return counts, nil
}
