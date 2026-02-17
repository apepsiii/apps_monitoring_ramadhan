package services

import (
	"log"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
	"github.com/ramadhan/amaliah-monitoring/internal/repository"
)

type BadgeService struct {
	BadgeRepo   *repository.BadgeRepository
	PrayerRepo  *repository.PrayerRepository
	AmaliahRepo *repository.AmaliahRepository
	QuranRepo   *repository.QuranRepository
}

func NewBadgeService(
	badgeRepo *repository.BadgeRepository,
	prayerRepo *repository.PrayerRepository,
	amaliahRepo *repository.AmaliahRepository,
	quranRepo *repository.QuranRepository,
) *BadgeService {
	return &BadgeService{
		BadgeRepo:   badgeRepo,
		PrayerRepo:  prayerRepo,
		AmaliahRepo: amaliahRepo,
		QuranRepo:   quranRepo,
	}
}

// CheckAndAwardBadges checks all criteria for a user and awards new badges
// Returns list of newly earned badges
func (s *BadgeService) CheckAndAwardBadges(userID int) ([]models.Badge, error) {
	var newBadges []models.Badge
	allBadges, err := s.BadgeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, badge := range allBadges {
		// Skip if already earned
		has, _ := s.BadgeRepo.HasBadge(userID, badge.ID)
		if has {
			continue
		}

		earned := false
		switch badge.CriteriaType {
		case "prayer_streak":
			streak, _, _ := s.PrayerRepo.GetPrayerStreak(userID)
			if streak >= badge.CriteriaValue {
				earned = true
			}
		case "amaliah_points":
			points, _ := s.AmaliahRepo.GetTotalPoints(userID, "2000-01-01", "2099-12-31") // implementation detail: verify GetTotalPoints
			if points >= badge.CriteriaValue {
				earned = true
			}
		case "quran_readings":
			count, _ := s.QuranRepo.GetTotalReadings(userID)
			if count >= badge.CriteriaValue {
				earned = true
			}
		// Add more criteria logic here
		}

		if earned {
			err := s.BadgeRepo.AwardBadge(userID, badge.ID)
			if err == nil {
				newBadges = append(newBadges, badge)
				log.Printf("User %d earned badge: %s", userID, badge.Name)
			}
		}
	}

	return newBadges, nil
}

// Specific check to avoid checking everything every time
func (s *BadgeService) CheckPrayerBadges(userID int) ([]models.Badge, error) {
	// ... logic similar to above but filtered for prayer related badges
	// keeping it simple for now and calling valid logic
	return s.CheckAndAwardBadges(userID) 
}
