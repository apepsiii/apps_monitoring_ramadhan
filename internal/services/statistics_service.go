package services

import (
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
	"github.com/ramadhan/amaliah-monitoring/internal/repository"
)

type StatisticsService struct {
	PrayerRepo  *repository.PrayerRepository
	AmaliahRepo *repository.AmaliahRepository
	FastingRepo *repository.FastingRepository
	UserRepo    *repository.UserRepository
}

func NewStatisticsService(
	prayerRepo *repository.PrayerRepository,
	amaliahRepo *repository.AmaliahRepository,
	fastingRepo *repository.FastingRepository,
	userRepo *repository.UserRepository,
) *StatisticsService {
	return &StatisticsService{
		PrayerRepo:  prayerRepo,
		AmaliahRepo: amaliahRepo,
		FastingRepo: fastingRepo,
		UserRepo:    userRepo,
	}
}

func (s *StatisticsService) GetDashboardStats() (map[string]interface{}, error) {
	// 1. Prayer Stats (Last 7 days)
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	
	prayerStats, err := s.PrayerRepo.GetDailyCompletionStats(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 2. Amaliah Distribution (Top 5)
	amaliahStats, err := s.AmaliahRepo.GetAmaliahDistribution()
	if err != nil {
		return nil, err
	}

	// 3. Fasting Stats (Today)
	// For simplicity, we can fetch all fasting records for today and count
	// Ideally should be a repo method, but we can reuse repo logic if available or just skip for now
	// Let's create a specialized query in repo if needed, or just return empty for MVP
	
	// Prepare dates for chart labels
	var labels []string
	var data []float64

	// Fill gaps if date is missing (optional, skipping for MVP)
	for _, stat := range prayerStats {
		dateStr := stat["date"].(string)
		parsed, _ := time.Parse("2006-01-02", dateStr)
		labels = append(labels, parsed.Format("02 Jan"))
		data = append(data, stat["percentage"].(float64))
	}

	// 4. Leaderboard (Top 10)
	topStudents, err := s.UserRepo.GetTopStudents(10)
	if err != nil {
		// Log error but don't fail the whole dashboard
		topStudents = []*models.User{}
	}

	return map[string]interface{}{
		"prayer_labels": labels,
		"prayer_data":   data,
		"amaliah_stats": amaliahStats,
		"top_students":  topStudents,
	}, nil
}
