package repository

import (
	"VoAr/internal/models"
	"context"
)

type CommentRepository struct {
	DB DBTX
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment models.Comment) error {
	_, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO comments (article_id, user_id, content, created_at) VALUES ($1, $2, $3, $4)",
		comment.ArticleID, comment.UserID, comment.Content, comment.CreatedAt,
	)
	return err
}

func (r *CommentRepository) GetCommentsByArticle(ctx context.Context, articleID int) ([]models.Comment, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, article_id, user_id, content, created_at, updated_at FROM comments WHERE article_id = $1",
		articleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
