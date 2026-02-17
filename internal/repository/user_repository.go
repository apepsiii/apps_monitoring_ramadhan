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
	query := `SELECT id, username, email, password_hash, full_name, class, role, points, avatar, bio, theme, target_khatam, created_at, updated_at
			  FROM users WHERE id = ?`

	user := &models.User{}
	err := r.DB.QueryRow(query, id).Scan(
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

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, class, role, points, avatar, bio, theme, target_khatam, created_at, updated_at
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
	query := `SELECT id, username, email, password_hash, full_name, class, role, points, avatar, bio, theme, target_khatam, created_at, updated_at
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
			  avatar = ?, theme = ?, target_khatam = ?, updated_at = ? 
			  WHERE id = ?`

	_, err := r.DB.Exec(query, req.FullName, req.Email, req.Class, req.Bio,
		req.Avatar, req.Theme, req.TargetKhatam, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdatePassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, hashedPassword, time.Now(), userID)
	return err
}
