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
