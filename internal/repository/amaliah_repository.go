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
