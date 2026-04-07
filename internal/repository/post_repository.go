package repository

import (
	"VoAr/internal/models"
	"database/sql"
)

type PostRepository struct {
	DB *sql.DB
}

func (r *PostRepository) GetPosts(page, pageSize int) ([]models.Post, error) {
	offset := (page - 1) * pageSize

	rows, err := r.DB.Query(
		"SELECT id, title, anons, full_text FROM articles LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
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
