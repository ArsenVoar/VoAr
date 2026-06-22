package repository

import (
	"VoAr/internal/models"
	"context"
	"database/sql"
	"errors"
)

type PostRepository struct {
	DB DBTX
}

func (r *PostRepository) GetPosts(ctx context.Context, page, pageSize int) ([]models.Post, error) {
	offset := (page - 1) * pageSize

	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, user_id, title, anons, full_text FROM articles LIMIT $1 OFFSET $2",
		pageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPostById(ctx context.Context, id int) (models.Post, error) {
	var post models.Post

	row := r.DB.QueryRowContext(
		ctx,
		"SELECT id, user_id, title, anons, full_text FROM articles WHERE id = $1",
		id,
	)

	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Anons, &post.FullText)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ErrNotFound
		}
		return models.Post{}, err
	}

	return post, nil
}
