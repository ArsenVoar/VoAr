package repository

import (
	"VoAr/internal/models"
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("not found")

type UserRepository struct {
	DB DBTX
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (models.User, error) {
	var user models.User

	row := r.DB.QueryRowContext(
		ctx,
		"SELECT id, name, email FROM users WHERE id = $1",
		id,
	)

	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	row := r.DB.QueryRowContext(
		ctx,
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	)

	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		user.Name, user.Email, user.Password,
	)
	return err
}
