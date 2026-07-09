package repository

import (
	"VoAr/internal/database"
	"VoAr/internal/models"
	"context"
)

type CommentRepository struct {
	DB database.DBTX
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment models.Comment) (int, error) {
	var commentID int

	err := r.DB.QueryRowContext(
		ctx,
		"INSERT INTO comments (post_id, user_id, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		comment.PostID, comment.UserID, comment.Content, comment.CreatedAt,
	).Scan(&commentID)
	return commentID, err
}

func (r *CommentRepository) GetCommentsByPost(ctx context.Context, postID int) ([]models.Comment, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, post_id, user_id, content, created_at, updated_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC, id ASC",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
