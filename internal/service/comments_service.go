package service

import (
	"VoAr/internal/logger"
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"errors"
	"time"
	"unicode/utf8"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment models.Comment) (int, error)
	GetCommentsByPost(ctx context.Context, postID int) ([]models.Comment, error)
}

type PostReader interface {
	GetPostByID(ctx context.Context, id int) (models.Post, error)
}

type NotificationNotifier interface {
	NotifyCommentCreated(ctx context.Context, recipientID int, actorID int, commentID int) error
}

type CommentService struct {
	CommentRepo CommentRepository
	PostReader  PostReader
	Notifier    NotificationNotifier
}

const maxCommentLength = 100

func (s *CommentService) CreateComment(ctx context.Context, comment models.Comment) (int, error) {
	if comment.Content == "" {
		return 0, ErrEmptyFields
	}
	if utf8.RuneCountInString(comment.Content) > maxCommentLength {
		return 0, ErrInvalidInput
	}

	post, err := s.PostReader.GetPostByID(ctx, comment.PostID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	comment.CreatedAt = time.Now()

	commentID, err := s.CommentRepo.CreateComment(ctx, comment)
	if err != nil {
		return 0, err
	}

	err = s.Notifier.NotifyCommentCreated(ctx, post.UserID, comment.UserID, commentID)
	if err != nil {
		logger.Error(err.Error())
	}

	return commentID, nil
}
func (s *CommentService) GetCommentsByPost(ctx context.Context, postID int) ([]models.Comment, error) {
	comments, err := s.CommentRepo.GetCommentsByPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
