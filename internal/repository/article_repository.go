package repository

import (
	"context"
	"database/sql"
	"errors"
)

type ArticleRepository struct {
	DB *sql.DB
}

func (r *ArticleRepository) Create(ctx context.Context, title, anons, fullText string) error {
	result, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO articles (title, anons, full_text) VALUES ($1, $2, $3)",
		title, anons, fullText,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
