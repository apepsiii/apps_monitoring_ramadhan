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

// Admin Methods

func (r *PrayerRepository) GetTodayStats(date string) (map[string]int, error) {
	query := `SELECT 
			  COUNT(DISTINCT user_id) as total_users,
			  COUNT(CASE WHEN subuh IN ('jamaah', 'sendiri') THEN 1 END) as subuh_count,
			  COUNT(CASE WHEN dzuhur IN ('jamaah', 'sendiri') THEN 1 END) as dzuhur_count,
			  COUNT(CASE WHEN ashar IN ('jamaah', 'sendiri') THEN 1 END) as ashar_count,
			  COUNT(CASE WHEN maghrib IN ('jamaah', 'sendiri') THEN 1 END) as maghrib_count,
			  COUNT(CASE WHEN isya IN ('jamaah', 'sendiri') THEN 1 END) as isya_count
			  FROM prayers 
			  WHERE date = ?`

	stats := make(map[string]int)
	var totalUsers, subuh, dzuhur, ashar, maghrib, isya int

	err := r.DB.QueryRow(query, date).Scan(&totalUsers, &subuh, &dzuhur, &ashar, &maghrib, &isya)
	if err != nil {
		return nil, err
	}

	stats["total_users"] = totalUsers
	stats["subuh"] = subuh
	stats["dzuhur"] = dzuhur
	stats["ashar"] = ashar
	stats["maghrib"] = maghrib
	stats["isya"] = isya

	return stats, nil
}

func (r *PrayerRepository) GetAllByDate(date string) ([]*models.Prayer, error) {
	query := `SELECT p.id, p.user_id, p.date, p.subuh, p.dzuhur, p.ashar, p.maghrib, p.isya, p.created_at, p.updated_at,
			  u.full_name, u.class
			  FROM prayers p
			  JOIN users u ON p.user_id = u.id
			  WHERE p.date = ?
			  ORDER BY u.class, u.full_name`

	rows, err := r.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prayers []*models.Prayer
	for rows.Next() {
		prayer := &models.Prayer{}
		var fullName, class string
		err := rows.Scan(
			&prayer.ID, &prayer.UserID, &prayer.Date, &prayer.Subuh,
			&prayer.Dzuhur, &prayer.Ashar, &prayer.Maghrib, &prayer.Isya,
			&prayer.CreatedAt, &prayer.UpdatedAt, &fullName, &class,
		)
		if err != nil {
			return nil, err
		}
		prayers = append(prayers, prayer)
	}
	return prayers, nil
}

func (r *PrayerRepository) GetByUser(userID int, limit int) ([]*models.Prayer, error) {
	query := `SELECT id, user_id, date, subuh, dzuhur, ashar, maghrib, isya, created_at, updated_at
			  FROM prayers WHERE user_id = ? ORDER BY date DESC LIMIT ?`

	rows, err := r.DB.Query(query, userID, limit)
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

func (r *PrayerRepository) GetPrayerStreak(userID int) (int, int, error) {
	query := `SELECT date, subuh, dzuhur, ashar, maghrib, isya 
			  FROM prayers 
			  WHERE user_id = ? 
			  ORDER BY date DESC 
			  LIMIT 60`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	type dayPrayer struct {
		date                                string
		subuh, dzuhur, ashar, maghrib, isya string
	}
	var prayers []dayPrayer
	for rows.Next() {
		var p dayPrayer
		err := rows.Scan(&p.date, &p.subuh, &p.dzuhur, &p.ashar, &p.maghrib, &p.isya)
		if err != nil {
			return 0, 0, err
		}
		prayers = append(prayers, p)
	}

	isComplete := func(p dayPrayer) bool {
		complete := 0
		if p.subuh == "jamaah" || p.subuh == "sendiri" {
			complete++
		}
		if p.dzuhur == "jamaah" || p.dzuhur == "sendiri" {
			complete++
		}
		if p.ashar == "jamaah" || p.ashar == "sendiri" {
			complete++
		}
		if p.maghrib == "jamaah" || p.maghrib == "sendiri" {
			complete++
		}
		if p.isya == "jamaah" || p.isya == "sendiri" {
			complete++
		}
		return complete == 5
	}

	currentStreak := 0
	bestStreak := 0
	tempStreak := 0

	for _, p := range prayers {
		if isComplete(p) {
			tempStreak++
			if tempStreak > bestStreak {
				bestStreak = tempStreak
			}
		} else {
			if currentStreak == 0 && len(prayers) > 0 {
				if isComplete(prayers[0]) {
					currentStreak = tempStreak
				}
			}
			tempStreak = 0
		}
	}

	if currentStreak == 0 && len(prayers) > 0 && isComplete(prayers[0]) {
		currentStreak = tempStreak
	}

	return currentStreak, bestStreak, nil
}
