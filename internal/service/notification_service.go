package service

import (
	"VoAr/internal/models"
	"context"
	"time"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification models.Notification) error
	GetNotificationsByUser(ctx context.Context, userID int) ([]models.Notification, error)
}

type NotificationService struct {
	NotificationRepo NotificationRepository
}

func (s *NotificationService) NotifyCommentCreated(ctx context.Context, recipientID int, actorID int, commentID int) error {
	if recipientID == actorID {
		return nil
	}

	notification := models.Notification{
		RecipientID: recipientID,
		ActorID:     actorID,
		CommentID:   commentID,
		CreatedAt:   time.Now(),
	}

	return s.NotificationRepo.CreateNotification(ctx, notification)
}

func (s *NotificationService) GetNotificationsByUser(ctx context.Context, userID int) ([]models.Notification, error) {
	return s.NotificationRepo.GetNotificationsByUser(ctx, userID)
}
