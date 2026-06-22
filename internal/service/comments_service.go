package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"errors"
	"log"
	"time"
)

type CommentService struct {
	CommentRepo         *repository.CommentRepository
	PostRepo            *repository.PostRepository
	NotificationService *NotificationService
}

const MaxCommentLength = 500

func (s *CommentService) CreateComment(ctx context.Context, comment models.Comment) (int, error) {
	if comment.Content == "" {
		return 0, ErrEmptyFields
	}

	post, err := s.PostRepo.GetPostById(ctx, comment.ArticleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	if len(comment.Content) > MaxCommentLength {
		return 0, ErrInvalidInput
	}

	comment.CreatedAt = time.Now()

	commentID, err := s.CommentRepo.CreateComment(ctx, comment)
	if err != nil {
		return 0, err
	}

	err = s.NotificationService.NotifyCommentCreated(ctx, post.UserID, comment.UserID, commentID)
	if err != nil {
		log.Printf("notification creation failed: %v", err)
	}

	return commentID, nil
}
func (s *CommentService) GetCommentsByArticle(ctx context.Context, articleID int) ([]models.Comment, error) {
	comments, err := s.CommentRepo.GetCommentsByArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
