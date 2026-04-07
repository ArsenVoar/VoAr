package repository

import (
	"VoAr/internal/models"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("not found")

type ShowPostRepository struct {
	DB *sql.DB
}

func (r *ShowPostRepository) GetPostById(id string) (models.Post, error) {
	var post models.Post

	row := r.DB.QueryRow(
		"SELECT id, title, anons, full_text FROM articles WHERE id = $1",
		id,
	)

	err := row.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ErrNotFound
		}
		return models.Post{}, err
	}

	return post, nil
}
