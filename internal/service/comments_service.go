package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"errors"
	"time"
)

type CommentService struct {
	CommentRepo *repository.CommentRepository
	PostRepo    *repository.PostRepository
}

const MaxCommentLength = 500

func (s *CommentService) CreateComment(ctx context.Context, comment models.Comment) error {
	if comment.Content == "" {
		return ErrEmptyFields
	}

	_, err := s.PostRepo.GetPostById(ctx, comment.ArticleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if len(comment.Content) > MaxCommentLength {
		return ErrInvalidInput
	}

	comment.CreatedAt = time.Now()

	err := s.CommentRepo.CreateComment(ctx, comment)
	if err != nil {
		return err
	}
	return nil
}
func (s *CommentService) GetCommentsByArticle(ctx context.Context, articleID int) ([]models.Comment, error) {
	comments, err := s.CommentRepo.GetCommentsByArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
