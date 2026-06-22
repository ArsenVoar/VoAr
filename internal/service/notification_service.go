package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"time"
)

type NotificationService struct {
	NotifyRepo *repository.NotificationRepository
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

	return s.NotifyRepo.CreateNotification(ctx, notification)
}

func (s *NotificationService) GetNotificationsByUser(ctx context.Context, userID int) ([]models.Notification, error) {
	return s.NotifyRepo.GetNotificationsByUser(ctx, userID)
}
