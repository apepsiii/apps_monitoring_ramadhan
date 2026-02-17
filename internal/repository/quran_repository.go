package repository

import (
	"database/sql"

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
