package repository

import (
	"VoAr/internal/database"
	"VoAr/internal/models"
	"context"
	"database/sql"
	"errors"
)

type PostRepository struct {
	DB database.DBTX
}

var ErrNoRowsAffected = errors.New("no rows affected")

func (r *PostRepository) CreatePost(ctx context.Context, userID int, title, anons, fullText string) error {
	result, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO posts (user_id, title, anons, full_text) VALUES ($1, $2, $3, $4)",
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
		return ErrNoRowsAffected
	}

	return nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, id int) (models.Post, error) {
	var post models.Post

	row := r.DB.QueryRowContext(
		ctx,
		"SELECT id, user_id, title, anons, full_text FROM posts WHERE id = $1",
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

func (r *PostRepository) GetPosts(ctx context.Context, page, pageSize int) ([]models.Post, error) {
	offset := (page - 1) * pageSize

	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, user_id, title, anons, full_text FROM posts ORDER BY id DESC LIMIT $1 OFFSET $2",
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
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Anons, &post.FullText); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
