package repository

import (
	"database/sql"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type PrayerRepository struct {
	DB *sql.DB
}

func NewPrayerRepository(db *sql.DB) *PrayerRepository {
	return &PrayerRepository{DB: db}
}

func (r *PrayerRepository) Create(prayer *models.Prayer) error {
	query := `INSERT INTO prayers (user_id, date, subuh, dzuhur, ashar, maghrib, isya) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.DB.Exec(query, prayer.UserID, prayer.Date, prayer.Subuh, 
		prayer.Dzuhur, prayer.Ashar, prayer.Maghrib, prayer.Isya)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	prayer.ID = int(id)
	return nil
}

func (r *PrayerRepository) GetByUserAndDate(userID int, date string) (*models.Prayer, error) {
	query := `SELECT id, user_id, date, subuh, dzuhur, ashar, maghrib, isya, created_at, updated_at 
			  FROM prayers WHERE user_id = ? AND date = ?`
	
	prayer := &models.Prayer{}
	err := r.DB.QueryRow(query, userID, date).Scan(
		&prayer.ID, &prayer.UserID, &prayer.Date, &prayer.Subuh,
		&prayer.Dzuhur, &prayer.Ashar, &prayer.Maghrib, &prayer.Isya,
		&prayer.CreatedAt, &prayer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return prayer, nil
}

func (r *PrayerRepository) Update(prayer *models.Prayer) error {
	query := `UPDATE prayers SET subuh = ?, dzuhur = ?, ashar = ?, maghrib = ?, isya = ?, updated_at = ? 
			  WHERE id = ?`
	
	_, err := r.DB.Exec(query, prayer.Subuh, prayer.Dzuhur, prayer.Ashar, 
		prayer.Maghrib, prayer.Isya, time.Now(), prayer.ID)
	return err
}

func (r *PrayerRepository) GetByUserAndDateRange(userID int, startDate, endDate string) ([]*models.Prayer, error) {
	query := `SELECT id, user_id, date, subuh, dzuhur, ashar, maghrib, isya, created_at, updated_at 
			  FROM prayers WHERE user_id = ? AND date BETWEEN ? AND ? ORDER BY date DESC`
	
	rows, err := r.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prayers []*models.Prayer
	for rows.Next() {
		prayer := &models.Prayer{}
		err := rows.Scan(
			&prayer.ID, &prayer.UserID, &prayer.Date, &prayer.Subuh,
			&prayer.Dzuhur, &prayer.Ashar, &prayer.Maghrib, &prayer.Isya,
			&prayer.CreatedAt, &prayer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		prayers = append(prayers, prayer)
	}
	return prayers, nil
}

func (r *PrayerRepository) GetTodayPrayer(userID int) (*models.Prayer, error) {
	today := time.Now().Format("2006-01-02")
	return r.GetByUserAndDate(userID, today)
}

func (r *PrayerRepository) CreateOrUpdate(userID int, date string, subuh, dzuhur, ashar, maghrib, isya string) error {
	existing, err := r.GetByUserAndDate(userID, date)
	if err != nil {
		// Create new
		prayer := &models.Prayer{
			UserID:  userID,
			Date:    date,
			Subuh:   subuh,
			Dzuhur:  dzuhur,
			Ashar:   ashar,
			Maghrib: maghrib,
			Isya:    isya,
		}
		return r.Create(prayer)
	}

	// Update existing
	existing.Subuh = subuh
	existing.Dzuhur = dzuhur
	existing.Ashar = ashar
	existing.Maghrib = maghrib
	existing.Isya = isya
	return r.Update(existing)
}

func (r *PrayerRepository) GetPrayerStats(userID int, startDate, endDate string) (map[string]int, error) {
	query := `SELECT 
			  COUNT(CASE WHEN subuh IN ('jamaah', 'sendiri') THEN 1 END) as subuh_count,
			  COUNT(CASE WHEN dzuhur IN ('jamaah', 'sendiri') THEN 1 END) as dzuhur_count,
			  COUNT(CASE WHEN ashar IN ('jamaah', 'sendiri') THEN 1 END) as ashar_count,
			  COUNT(CASE WHEN maghrib IN ('jamaah', 'sendiri') THEN 1 END) as maghrib_count,
			  COUNT(CASE WHEN isya IN ('jamaah', 'sendiri') THEN 1 END) as isya_count,
			  COUNT(*) as total_days
			  FROM prayers 
			  WHERE user_id = ? AND date BETWEEN ? AND ?`
	
	stats := make(map[string]int)
	var subuh, dzuhur, ashar, maghrib, isya, total int
	
	err := r.DB.QueryRow(query, userID, startDate, endDate).Scan(&subuh, &dzuhur, &ashar, &maghrib, &isya, &total)
	if err != nil {
		return nil, err
	}

	stats["subuh"] = subuh
	stats["dzuhur"] = dzuhur
	stats["ashar"] = ashar
	stats["maghrib"] = maghrib
	stats["isya"] = isya
	stats["total_days"] = total

	return stats, nil
}
