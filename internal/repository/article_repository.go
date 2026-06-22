package repository

import (
	"context"
	"errors"
)

type ArticleRepository struct {
	DB DBTX
}

func (r *ArticleRepository) Create(ctx context.Context, userID int, title, anons, fullText string) error {
	result, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO articles (user_id, title, anons, full_text) VALUES ($1, $2, $3, $4)",
		userID, title, anons, fullText,
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
