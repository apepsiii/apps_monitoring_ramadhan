package repository

import (
	"database/sql"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type QuranRepository struct {
	DB *sql.DB
}

func NewQuranRepository(db *sql.DB) *QuranRepository {
	return &QuranRepository{DB: db}
}

func (r *QuranRepository) Create(reading *models.QuranReading) error {
	query := `INSERT INTO quran_readings (user_id, date, start_surah_id, start_surah_name, start_ayah, end_surah_id, end_surah_name, end_ayah, notes)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query, reading.UserID, reading.Date, reading.StartSurahID,
		reading.StartSurahName, reading.StartAyah, reading.EndSurahID, reading.EndSurahName, reading.EndAyah, reading.Notes)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	reading.ID = int(id)
	return nil
}

func (r *QuranRepository) GetByUserAndDate(userID int, date string) ([]*models.QuranReading, error) {
	query := `SELECT id, user_id, date, start_surah_id, start_surah_name, start_ayah, end_surah_id, end_surah_name, end_ayah, notes, created_at
			  FROM quran_readings WHERE user_id = ? AND date = ? ORDER BY created_at DESC`

	rows, err := r.DB.Query(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*models.QuranReading
	for rows.Next() {
		reading := &models.QuranReading{}
		err := rows.Scan(
			&reading.ID, &reading.UserID, &reading.Date, &reading.StartSurahID, &reading.StartSurahName,
			&reading.StartAyah, &reading.EndSurahID, &reading.EndSurahName, &reading.EndAyah, &reading.Notes, &reading.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, nil
}

func (r *QuranRepository) GetByUser(userID int, limit int) ([]*models.QuranReading, error) {
	query := `SELECT id, user_id, date, start_surah_id, start_surah_name, start_ayah, end_surah_id, end_surah_name, end_ayah, notes, created_at
			  FROM quran_readings WHERE user_id = ? ORDER BY date DESC, created_at DESC LIMIT ?`

	rows, err := r.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*models.QuranReading
	for rows.Next() {
		reading := &models.QuranReading{}
		err := rows.Scan(
			&reading.ID, &reading.UserID, &reading.Date, &reading.StartSurahID, &reading.StartSurahName,
			&reading.StartAyah, &reading.EndSurahID, &reading.EndSurahName, &reading.EndAyah, &reading.Notes, &reading.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, nil
}

func (r *QuranRepository) GetTotalReadings(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM quran_readings WHERE user_id = ?`

	var total int
	err := r.DB.QueryRow(query, userID).Scan(&total)
	return total, err
}

func (r *QuranRepository) GetTotalPagesRead(userID int, startDate, endDate string) (int, error) {
	query := `SELECT COALESCE(SUM(pages), 0) FROM quran_readings 
			  WHERE user_id = ? AND date BETWEEN ? AND ?`

	var total int
	err := r.DB.QueryRow(query, userID, startDate, endDate).Scan(&total)
	return total, err
}

func (r *QuranRepository) Delete(id int) error {
	query := `DELETE FROM quran_readings WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

// Admin Methods

func (r *QuranRepository) GetTodayStats(date string) (map[string]interface{}, error) {
	query := `SELECT 
			  COUNT(DISTINCT user_id) as total_users,
			  COUNT(*) as total_readings
			  FROM quran_readings 
			  WHERE date = ?`

	stats := make(map[string]interface{})
	var totalUsers, totalReadings int

	err := r.DB.QueryRow(query, date).Scan(&totalUsers, &totalReadings)
	if err != nil {
		return nil, err
	}

	stats["total_users"] = totalUsers
	stats["total_readings"] = totalReadings

	return stats, nil
}

func (r *QuranRepository) GetAllByDate(date string) ([]*models.QuranReading, error) {
	query := `SELECT qr.id, qr.user_id, qr.date, qr.start_surah_id, qr.start_surah_name, qr.start_ayah,
			  qr.end_surah_id, qr.end_surah_name, qr.end_ayah, qr.notes, qr.created_at,
			  u.full_name, u.class
			  FROM quran_readings qr
			  JOIN users u ON qr.user_id = u.id
			  WHERE qr.date = ?
			  ORDER BY u.class, u.full_name`

	rows, err := r.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []*models.QuranReading
	for rows.Next() {
		reading := &models.QuranReading{}
		var fullName, class string
		err := rows.Scan(
			&reading.ID, &reading.UserID, &reading.Date, &reading.StartSurahID, &reading.StartSurahName,
			&reading.StartAyah, &reading.EndSurahID, &reading.EndSurahName, &reading.EndAyah,
			&reading.Notes, &reading.CreatedAt, &fullName, &class,
		)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, nil
}

func (r *QuranRepository) GetQuranStreak(userID int) (int, int, error) {
	query := `SELECT DISTINCT date FROM quran_readings 
			  WHERE user_id = ? 
			  ORDER BY date DESC 
			  LIMIT 60`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		err := rows.Scan(&date)
		if err != nil {
			return 0, 0, err
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0, 0, nil
	}

	parseDate := func(s string) int {
		t, _ := time.Parse("2006-01-02", s)
		return int(t.Unix() / 86400)
	}

	currentStreak := 0
	bestStreak := 0
	tempStreak := 1

	for i := 1; i < len(dates); i++ {
		diff := parseDate(dates[i-1]) - parseDate(dates[i])
		if diff == 1 {
			tempStreak++
		} else {
			if bestStreak < tempStreak {
				bestStreak = tempStreak
			}
			if currentStreak == 0 {
				currentStreak = tempStreak
			}
			tempStreak = 1
		}
	}

	if bestStreak < tempStreak {
		bestStreak = tempStreak
	}
	if currentStreak == 0 {
		currentStreak = tempStreak
	}

	today := time.Now().Format("2006-01-02")
	if dates[0] != today {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if dates[0] != yesterday {
			return 0, bestStreak, nil
		}
	}

	return currentStreak, bestStreak, nil
}
