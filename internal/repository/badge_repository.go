package repository

import (
	"database/sql"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type BadgeRepository struct {
	DB *sql.DB
}

func NewBadgeRepository(db *sql.DB) *BadgeRepository {
	return &BadgeRepository{DB: db}
}

func (r *BadgeRepository) GetAll() ([]models.Badge, error) {
	query := `SELECT id, name, description, icon, criteria_type, criteria_value, created_at FROM badges`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []models.Badge
	for rows.Next() {
		var b models.Badge
		err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Icon, &b.CriteriaType, &b.CriteriaValue, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *BadgeRepository) GetUserBadges(userID int) ([]models.UserBadge, error) {
	query := `SELECT ub.id, ub.user_id, ub.badge_id, ub.earned_at,
			  b.id, b.name, b.description, b.icon, b.criteria_type, b.criteria_value
			  FROM user_badges ub
			  JOIN badges b ON ub.badge_id = b.id
			  WHERE ub.user_id = ?
			  ORDER BY ub.earned_at DESC`
	
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBadges []models.UserBadge
	for rows.Next() {
		var ub models.UserBadge
		var b models.Badge
		err := rows.Scan(&ub.ID, &ub.UserID, &ub.BadgeID, &ub.EarnedAt,
			&b.ID, &b.Name, &b.Description, &b.Icon, &b.CriteriaType, &b.CriteriaValue)
		if err != nil {
			return nil, err
		}
		ub.Badge = b
		userBadges = append(userBadges, ub)
	}
	return userBadges, nil
}

func (r *BadgeRepository) HasBadge(userID, badgeID int) (bool, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM user_badges WHERE user_id = ? AND badge_id = ?", userID, badgeID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BadgeRepository) AwardBadge(userID, badgeID int) error {
	query := `INSERT INTO user_badges (user_id, badge_id, earned_at) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, userID, badgeID, time.Now())
	return err
}

func (r *BadgeRepository) GetbadgesByCriteria(criteriaType string) ([]models.Badge, error) {
	query := `SELECT id, name, description, icon, criteria_type, criteria_value, created_at FROM badges WHERE criteria_type = ?`
	rows, err := r.DB.Query(query, criteriaType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []models.Badge
	for rows.Next() {
		var b models.Badge
		err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Icon, &b.CriteriaType, &b.CriteriaValue, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}
