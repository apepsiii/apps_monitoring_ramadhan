package repository

import (
	"database/sql"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type ClassRepository struct {
	DB *sql.DB
}

func NewClassRepository(db *sql.DB) *ClassRepository {
	return &ClassRepository{DB: db}
}

func (r *ClassRepository) GetAll() ([]*models.Class, error) {
	query := `SELECT id, name, level, description, created_at, updated_at 
			  FROM classes ORDER BY name`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*models.Class
	for rows.Next() {
		c := &models.Class{}
		err := rows.Scan(
			&c.ID, &c.Name, &c.Level, &c.Description,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, nil
}

func (r *ClassRepository) Create(class *models.Class) error {
	query := `INSERT INTO classes (name, level, description) VALUES (?, ?, ?)`

	result, err := r.DB.Exec(query, class.Name, class.Level, class.Description)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	class.ID = int(id)
	return nil
}

func (r *ClassRepository) Update(class *models.Class) error {
	query := `UPDATE classes SET name = ?, level = ?, description = ?, updated_at = ? 
			  WHERE id = ?`

	_, err := r.DB.Exec(query, class.Name, class.Level, class.Description, time.Now(), class.ID)
	return err
}

func (r *ClassRepository) Delete(id int) error {
	query := `DELETE FROM classes WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *ClassRepository) GetByID(id int) (*models.Class, error) {
	query := `SELECT id, name, level, description, created_at, updated_at 
			  FROM classes WHERE id = ?`

	c := &models.Class{}
	err := r.DB.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Level, &c.Description,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
