package repository

import (
	"VoAr/internal/models"
	"database/sql"
	"errors"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetByID(id string) (models.User, error) {
	var user models.User

	row := r.DB.QueryRow(
		"SELECT id, name, email FROM users WHERE id = $1",
		id,
	)

	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

var ErrUserExists = errors.New("user already exists")

func (r *UserRepository) Create(user models.User) error {
	_, err := r.DB.Exec(
		"INSERT INTO users (name, email) VALUES ($1, $2)",
		user.Name, user.Email,
	)
	if err != nil {
		// обработка pq ошибки
		return err
	}
	return nil
}
