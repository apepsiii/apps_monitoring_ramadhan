package repository

import (
	"database/sql"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, full_name, class, role, points, bio) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query, user.Username, user.Email, user.PasswordHash,
		user.FullName, user.Class, user.Role, user.Points, "")
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, 
			  COALESCE(class, '') as class, role, points, 
			  COALESCE(avatar, 'default') as avatar, 
			  COALESCE(bio, '') as bio, 
			  COALESCE(theme, 'emerald') as theme, 
			  COALESCE(target_khatam, 30) as target_khatam, 
			  COALESCE(provinsi, '') as provinsi,
			  COALESCE(kabkota, '') as kabkota,
			  created_at, updated_at
			  FROM users WHERE id = ?`

	user := &models.User{}
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.Class, &user.Role, &user.Points,
		&user.Avatar, &user.Bio, &user.Theme, &user.TargetKhatam,
		&user.Provinsi, &user.Kabkota,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, 
			  COALESCE(class, '') as class, role, points, 
			  COALESCE(avatar, 'default') as avatar, 
			  COALESCE(bio, '') as bio, 
			  COALESCE(theme, 'emerald') as theme, 
			  COALESCE(target_khatam, 30) as target_khatam, 
			  created_at, updated_at
			  FROM users WHERE username = ?`

	user := &models.User{}
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.Class, &user.Role, &user.Points,
		&user.Avatar, &user.Bio, &user.Theme, &user.TargetKhatam,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, 
			  COALESCE(class, '') as class, role, points, 
			  COALESCE(avatar, 'default') as avatar, 
			  COALESCE(bio, '') as bio, 
			  COALESCE(theme, 'emerald') as theme, 
			  COALESCE(target_khatam, 30) as target_khatam, 
			  created_at, updated_at
			  FROM users WHERE email = ?`

	user := &models.User{}
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.Class, &user.Role, &user.Points,
		&user.Avatar, &user.Bio, &user.Theme, &user.TargetKhatam,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET username = ?, email = ?, full_name = ?, class = ?, 
			  points = ?, updated_at = ? WHERE id = ?`

	_, err := r.DB.Exec(query, user.Username, user.Email, user.FullName,
		user.Class, user.Points, time.Now(), user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := `SELECT id, username, email, full_name, class, role, points, avatar, bio, theme, target_khatam, created_at, updated_at
			  FROM users ORDER BY created_at DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FullName,
			&user.Class, &user.Role, &user.Points, &user.Avatar, &user.Bio,
			&user.Theme, &user.TargetKhatam, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) UpdatePoints(userID int, points int) error {
	query := `UPDATE users SET points = points + ? WHERE id = ?`
	_, err := r.DB.Exec(query, points, userID)
	return err
}

func (r *UserRepository) UpdateProfile(userID int, req *models.ProfileUpdateRequest) error {
	query := `UPDATE users SET 
			  full_name = ?, email = ?, class = ?, bio = ?, 
			  avatar = ?, theme = ?, target_khatam = ?, 
			  provinsi = ?, kabkota = ?, updated_at = ? 
			  WHERE id = ?`

	_, err := r.DB.Exec(query, req.FullName, req.Email, req.Class, req.Bio,
		req.Avatar, req.Theme, req.TargetKhatam,
		req.Provinsi, req.Kabkota, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdatePassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, hashedPassword, time.Now(), userID)
	return err
}

// Admin Methods

func (r *UserRepository) GetStats() (map[string]interface{}, error) {
	query := `SELECT 
			  COUNT(*) as total_users,
			  COUNT(CASE WHEN role = 'admin' THEN 1 END) as admin_count,
			  COUNT(CASE WHEN role = 'user' THEN 1 END) as user_count,
			  COALESCE(SUM(points), 0) as total_points,
			  COUNT(DISTINCT class) as total_classes
			  FROM users`

	stats := make(map[string]interface{})
	var totalUsers, adminCount, userCount, totalPoints, totalClasses int

	err := r.DB.QueryRow(query).Scan(&totalUsers, &adminCount, &userCount, &totalPoints, &totalClasses)
	if err != nil {
		return nil, err
	}

	stats["total_users"] = totalUsers
	stats["admin_count"] = adminCount
	stats["user_count"] = userCount
	stats["total_points"] = totalPoints
	stats["total_classes"] = totalClasses

	return stats, nil
}

func (r *UserRepository) GetActiveUsersCount(date string) (int, error) {
	query := `SELECT COUNT(DISTINCT user_id) FROM (
			  SELECT user_id FROM prayers WHERE date = ?
			  UNION
			  SELECT user_id FROM fastings WHERE date = ?
			  UNION
			  SELECT user_id FROM quran_readings WHERE date = ?
			  UNION
			  SELECT user_id FROM daily_amaliah WHERE date = ?
			)`

	var count int
	err := r.DB.QueryRow(query, date, date, date, date).Scan(&count)
	return count, err
}

func (r *UserRepository) GetByClass(class string) ([]*models.User, error) {
	query := `SELECT id, username, email, full_name, class, role, points, avatar, bio, theme, target_khatam, created_at, updated_at
			  FROM users WHERE class = ? AND role = 'user' ORDER BY full_name`

	rows, err := r.DB.Query(query, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FullName,
			&user.Class, &user.Role, &user.Points, &user.Avatar, &user.Bio,
			&user.Theme, &user.TargetKhatam, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetAllClasses() ([]string, error) {
	query := `SELECT DISTINCT class FROM users WHERE role = 'user' AND class IS NOT NULL AND class != '' ORDER BY class`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []string
	for rows.Next() {
		var class string
		err := rows.Scan(&class)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (r *UserRepository) SearchUsers(query string) ([]*models.User, error) {
	sqlQuery := `SELECT id, username, email, full_name, class, role, points, avatar, bio, theme, target_khatam, created_at, updated_at
				 FROM users 
				 WHERE role = 'user' AND 
				 (full_name LIKE ? OR username LIKE ? OR email LIKE ? OR class LIKE ?)
				 ORDER BY full_name`

	searchTerm := "%" + query + "%"
	rows, err := r.DB.Query(sqlQuery, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FullName,
			&user.Class, &user.Role, &user.Points, &user.Avatar, &user.Bio,
			&user.Theme, &user.TargetKhatam, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
