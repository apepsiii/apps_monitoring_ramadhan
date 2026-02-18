package repository

import (
	"database/sql"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type FastingRepository struct {
	DB *sql.DB
}

func NewFastingRepository(db *sql.DB) *FastingRepository {
	return &FastingRepository{DB: db}
}

func (r *FastingRepository) Create(fasting *models.Fasting) error {
	query := `INSERT INTO fastings (user_id, date, status, reason) VALUES (?, ?, ?, ?)`

	result, err := r.DB.Exec(query, fasting.UserID, fasting.Date, fasting.Status, fasting.Reason)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	fasting.ID = int(id)
	return nil
}

func (r *FastingRepository) GetByUserAndDate(userID int, date string) (*models.Fasting, error) {
	query := `SELECT id, user_id, date, status, reason, created_at 
			  FROM fastings WHERE user_id = ? AND date = ?`

	fasting := &models.Fasting{}
	err := r.DB.QueryRow(query, userID, date).Scan(
		&fasting.ID, &fasting.UserID, &fasting.Date, &fasting.Status,
		&fasting.Reason, &fasting.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return fasting, nil
}

func (r *FastingRepository) Update(fasting *models.Fasting) error {
	query := `UPDATE fastings SET status = ?, reason = ? WHERE id = ?`

	_, err := r.DB.Exec(query, fasting.Status, fasting.Reason, fasting.ID)
	return err
}

func (r *FastingRepository) GetByUserAndDateRange(userID int, startDate, endDate string) ([]*models.Fasting, error) {
	query := `SELECT id, user_id, date, status, reason, created_at 
			  FROM fastings WHERE user_id = ? AND date BETWEEN ? AND ? ORDER BY date DESC`

	rows, err := r.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fastings []*models.Fasting
	for rows.Next() {
		fasting := &models.Fasting{}
		err := rows.Scan(
			&fasting.ID, &fasting.UserID, &fasting.Date, &fasting.Status,
			&fasting.Reason, &fasting.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		fastings = append(fastings, fasting)
	}
	return fastings, nil
}

func (r *FastingRepository) GetTodayFasting(userID int) (*models.Fasting, error) {
	today := time.Now().Format("2006-01-02")
	return r.GetByUserAndDate(userID, today)
}

func (r *FastingRepository) CreateOrUpdate(userID int, date, status, reason string) error {
	existing, err := r.GetByUserAndDate(userID, date)
	if err != nil {
		// Create new
		fasting := &models.Fasting{
			UserID: userID,
			Date:   date,
			Status: status,
			Reason: reason,
		}
		return r.Create(fasting)
	}

	// Update existing
	existing.Status = status
	existing.Reason = reason
	return r.Update(existing)
}

func (r *FastingRepository) GetFastingStats(userID int, startDate, endDate string) (map[string]int, error) {
	query := `SELECT 
			  COUNT(CASE WHEN status = 'puasa' THEN 1 END) as fasting_count,
			  COUNT(CASE WHEN status = 'tidak' THEN 1 END) as not_fasting_count,
			  COUNT(*) as total_days
			  FROM fastings 
			  WHERE user_id = ? AND date BETWEEN ? AND ?`

	stats := make(map[string]int)
	var fasting, notFasting, total int

	err := r.DB.QueryRow(query, userID, startDate, endDate).Scan(&fasting, &notFasting, &total)
	if err != nil {
		return nil, err
	}

	stats["fasting"] = fasting
	stats["not_fasting"] = notFasting
	stats["total_days"] = total

	return stats, nil
}

func (r *FastingRepository) GetTotalFasting(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM fastings WHERE user_id = ? AND status = 'puasa'`
	var total int
	err := r.DB.QueryRow(query, userID).Scan(&total)
	return total, err
}

// Admin Methods

func (r *FastingRepository) GetTodayStats(date string) (map[string]int, error) {
	query := `SELECT 
			  COUNT(DISTINCT user_id) as total_users,
			  COUNT(CASE WHEN status = 'puasa' THEN 1 END) as fasting_count,
			  COUNT(CASE WHEN status = 'tidak' THEN 1 END) as not_fasting_count
			  FROM fastings 
			  WHERE date = ?`

	stats := make(map[string]int)
	var totalUsers, fasting, notFasting int

	err := r.DB.QueryRow(query, date).Scan(&totalUsers, &fasting, &notFasting)
	if err != nil {
		return nil, err
	}

	stats["total_users"] = totalUsers
	stats["fasting"] = fasting
	stats["not_fasting"] = notFasting

	return stats, nil
}

func (r *FastingRepository) GetAllByDate(date string) ([]*models.Fasting, error) {
	query := `SELECT f.id, f.user_id, f.date, f.status, f.reason, f.created_at,
			  u.full_name, u.class
			  FROM fastings f
			  JOIN users u ON f.user_id = u.id
			  WHERE f.date = ?
			  ORDER BY u.class, u.full_name`

	rows, err := r.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fastings []*models.Fasting
	for rows.Next() {
		fasting := &models.Fasting{}
		var fullName, class string
		err := rows.Scan(
			&fasting.ID, &fasting.UserID, &fasting.Date, &fasting.Status,
			&fasting.Reason, &fasting.CreatedAt, &fullName, &class,
		)
		if err != nil {
			return nil, err
		}
		fastings = append(fastings, fasting)
	}
	return fastings, nil
}

func (r *FastingRepository) GetByUser(userID int, limit int) ([]*models.Fasting, error) {
	query := `SELECT id, user_id, date, status, reason, created_at
			  FROM fastings WHERE user_id = ? ORDER BY date DESC LIMIT ?`

	rows, err := r.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fastings []*models.Fasting
	for rows.Next() {
		fasting := &models.Fasting{}
		err := rows.Scan(
			&fasting.ID, &fasting.UserID, &fasting.Date, &fasting.Status,
			&fasting.Reason, &fasting.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		fastings = append(fastings, fasting)
	}
	return fastings, nil
}

func (r *FastingRepository) GetFastingStreak(userID int) (int, int, error) {
	query := `SELECT date, status FROM fastings 
			  WHERE user_id = ? 
			  ORDER BY date DESC 
			  LIMIT 60`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	type dayFasting struct {
		date   string
		status string
	}
	var fastings []dayFasting
	for rows.Next() {
		var f dayFasting
		err := rows.Scan(&f.date, &f.status)
		if err != nil {
			return 0, 0, err
		}
		fastings = append(fastings, f)
	}

	currentStreak := 0
	bestStreak := 0
	tempStreak := 0

	for i, f := range fastings {
		if f.status == "puasa" {
			tempStreak++
			if tempStreak > bestStreak {
				bestStreak = tempStreak
			}
		} else {
			if currentStreak == 0 && i > 0 {
				if len(fastings) > 0 && fastings[0].status == "puasa" {
					currentStreak = tempStreak
				}
			}
			tempStreak = 0
		}
	}

	if currentStreak == 0 && len(fastings) > 0 && fastings[0].status == "puasa" {
		currentStreak = tempStreak
	}

	return currentStreak, bestStreak, nil
}
