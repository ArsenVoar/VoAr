package models

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrUserExists = errors.New("user already exists")

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user User) error {
	_, err := r.DB.Exec(
		"INSERT INTO users (name, email) VALUES ($1, $2)",
		user.Name, user.Email,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrUserExists
		}
		return err
	}
	return nil
}
