package repository

import (
	"database/sql"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type AmaliahRepository struct {
	DB *sql.DB
}

func NewAmaliahRepository(db *sql.DB) *AmaliahRepository {
	return &AmaliahRepository{DB: db}
}

func (r *AmaliahRepository) GetAllTypes() ([]*models.AmaliahType, error) {
	query := `SELECT id, name, description, points, icon, is_active, created_at 
			  FROM amaliah_types WHERE is_active = 1 ORDER BY name`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*models.AmaliahType
	for rows.Next() {
		at := &models.AmaliahType{}
		err := rows.Scan(
			&at.ID, &at.Name, &at.Description, &at.Points, &at.Icon,
			&at.IsActive, &at.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		types = append(types, at)
	}
	return types, nil
}

func (r *AmaliahRepository) GetTypeByID(id int) (*models.AmaliahType, error) {
	query := `SELECT id, name, description, points, icon, is_active, created_at 
			  FROM amaliah_types WHERE id = ?`

	at := &models.AmaliahType{}
	err := r.DB.QueryRow(query, id).Scan(
		&at.ID, &at.Name, &at.Description, &at.Points, &at.Icon,
		&at.IsActive, &at.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (r *AmaliahRepository) CreateDailyAmaliah(da *models.DailyAmaliah) error {
	query := `INSERT INTO daily_amaliah (user_id, amaliah_type_id, date, notes) 
			  VALUES (?, ?, ?, ?)`

	result, err := r.DB.Exec(query, da.UserID, da.AmaliahTypeID, da.Date, da.Notes)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	da.ID = int(id)
	return nil
}

func (r *AmaliahRepository) GetDailyAmaliah(userID int, date string) ([]*models.DailyAmaliah, error) {
	query := `SELECT da.id, da.user_id, da.amaliah_type_id, da.date, da.notes, da.created_at,
			  at.id, at.name, at.description, at.points, at.icon
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  WHERE da.user_id = ? AND da.date = ?
			  ORDER BY da.created_at DESC`

	rows, err := r.DB.Query(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.DailyAmaliah
	for rows.Next() {
		item := &models.DailyAmaliah{}
		at := &models.AmaliahType{}
		err := rows.Scan(
			&item.ID, &item.UserID, &item.AmaliahTypeID, &item.Date, &item.Notes, &item.CreatedAt,
			&at.ID, &at.Name, &at.Description, &at.Points, &at.Icon,
		)
		if err != nil {
			return nil, err
		}
		item.AmaliahType = *at
		items = append(items, item)
	}
	return items, nil
}

func (r *AmaliahRepository) GetDailyAmaliahByType(userID int, amaliahTypeID int, date string) (*models.DailyAmaliah, error) {
	query := `SELECT id, user_id, amaliah_type_id, date, notes, created_at 
			  FROM daily_amaliah WHERE user_id = ? AND amaliah_type_id = ? AND date = ?`

	item := &models.DailyAmaliah{}
	err := r.DB.QueryRow(query, userID, amaliahTypeID, date).Scan(
		&item.ID, &item.UserID, &item.AmaliahTypeID, &item.Date, &item.Notes, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *AmaliahRepository) DeleteDailyAmaliah(id int) error {
	query := `DELETE FROM daily_amaliah WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *AmaliahRepository) GetTodayPoints(userID int) (int, error) {
	today := time.Now().Format("2006-01-02")
	query := `SELECT COALESCE(SUM(at.points), 0) 
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  WHERE da.user_id = ? AND da.date = ?`

	var points int
	err := r.DB.QueryRow(query, userID, today).Scan(&points)
	return points, err
}

func (r *AmaliahRepository) GetTotalPoints(userID int, startDate, endDate string) (int, error) {
	query := `SELECT COALESCE(SUM(at.points), 0) 
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  WHERE da.user_id = ? AND da.date BETWEEN ? AND ?`

	var points int
	err := r.DB.QueryRow(query, userID, startDate, endDate).Scan(&points)
	return points, err
}

func (r *AmaliahRepository) GetLeaderboard(limit int) ([]map[string]interface{}, error) {
	query := `SELECT u.id, u.full_name, u.class, u.points,
			  COUNT(DISTINCT da.date) as active_days
			  FROM users u
			  LEFT JOIN daily_amaliah da ON u.id = da.user_id
			  WHERE u.role = 'user'
			  GROUP BY u.id
			  ORDER BY u.points DESC
			  LIMIT ?`

	rows, err := r.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []map[string]interface{}
	for rows.Next() {
		var id, points, activeDays int
		var fullName, class string
		err := rows.Scan(&id, &fullName, &class, &points, &activeDays)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, map[string]interface{}{
			"id":          id,
			"full_name":   fullName,
			"class":       class,
			"points":      points,
			"active_days": activeDays,
		})
	}
	return leaderboard, nil
}

// Admin Methods

func (r *AmaliahRepository) GetTodayStats(date string) (map[string]interface{}, error) {
	query := `SELECT 
			  COUNT(DISTINCT da.user_id) as total_users,
			  COUNT(da.id) as total_amaliah,
			  COALESCE(SUM(at.points), 0) as total_points
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  WHERE da.date = ?`

	stats := make(map[string]interface{})
	var totalUsers, totalAmaliah, totalPoints int

	err := r.DB.QueryRow(query, date).Scan(&totalUsers, &totalAmaliah, &totalPoints)
	if err != nil {
		return nil, err
	}

	stats["total_users"] = totalUsers
	stats["total_amaliah"] = totalAmaliah
	stats["total_points"] = totalPoints

	return stats, nil
}

func (r *AmaliahRepository) GetStatsByType(date string) ([]map[string]interface{}, error) {
	query := `SELECT 
			  at.name, at.icon, COUNT(da.id) as count, COALESCE(SUM(at.points), 0) as points
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  WHERE da.date = ?
			  GROUP BY at.id
			  ORDER BY points DESC`

	rows, err := r.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []map[string]interface{}
	for rows.Next() {
		var name, icon string
		var count, points int
		err := rows.Scan(&name, &icon, &count, &points)
		if err != nil {
			return nil, err
		}
		stats = append(stats, map[string]interface{}{
			"name":   name,
			"icon":   icon,
			"count":  count,
			"points": points,
		})
	}
	return stats, nil
}

func (r *AmaliahRepository) GetAmaliahDistribution() ([]map[string]interface{}, error) {
	query := `SELECT 
			  at.name, COUNT(da.id) as count
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  GROUP BY at.id
			  ORDER BY count DESC
			  LIMIT 5`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []map[string]interface{}
	for rows.Next() {
		var name string
		var count int
		err := rows.Scan(&name, &count)
		if err != nil {
			return nil, err
		}
		stats = append(stats, map[string]interface{}{
			"name":  name,
			"count": count,
		})
	}
	return stats, nil
}

func (r *AmaliahRepository) GetAllByDate(date string) ([]*models.DailyAmaliah, error) {
	query := `SELECT da.id, da.user_id, da.amaliah_type_id, da.date, da.notes, da.created_at,
			  at.id, at.name, at.description, at.points, at.icon,
			  u.full_name, u.class
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  JOIN users u ON da.user_id = u.id
			  WHERE da.date = ?
			  ORDER BY u.class, u.full_name`

	rows, err := r.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.DailyAmaliah
	for rows.Next() {
		item := &models.DailyAmaliah{}
		at := &models.AmaliahType{}
		var fullName, class string
		err := rows.Scan(
			&item.ID, &item.UserID, &item.AmaliahTypeID, &item.Date, &item.Notes, &item.CreatedAt,
			&at.ID, &at.Name, &at.Description, &at.Points, &at.Icon,
			&fullName, &class,
		)
		if err != nil {
			return nil, err
		}
		item.AmaliahType = *at
		items = append(items, item)
	}
	return items, nil
}

func (r *AmaliahRepository) GetByUser(userID int, limit int) ([]*models.DailyAmaliah, error) {
	query := `SELECT da.id, da.user_id, da.amaliah_type_id, da.date, da.notes, da.created_at,
			  at.id, at.name, at.description, at.points, at.icon
			  FROM daily_amaliah da
			  JOIN amaliah_types at ON da.amaliah_type_id = at.id
			  WHERE da.user_id = ?
			  ORDER BY da.date DESC, da.created_at DESC
			  LIMIT ?`

	rows, err := r.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.DailyAmaliah
	for rows.Next() {
		item := &models.DailyAmaliah{}
		at := &models.AmaliahType{}
		err := rows.Scan(
			&item.ID, &item.UserID, &item.AmaliahTypeID, &item.Date, &item.Notes, &item.CreatedAt,
			&at.ID, &at.Name, &at.Description, &at.Points, &at.Icon,
		)
		if err != nil {
			return nil, err
		}
		item.AmaliahType = *at
		items = append(items, item)
	}
	return items, nil
}

func (r *AmaliahRepository) GetAmaliahStreak(userID int) (int, int, error) {
	query := `SELECT DISTINCT date FROM daily_amaliah 
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
