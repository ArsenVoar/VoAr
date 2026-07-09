package service

import (
	"VoAr/internal/models"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestNotifyCommentCreated_RecipientIsActor(t *testing.T) {
	notification := newTestNotification()
	notification.ActorID = notification.RecipientID

	fakeNotificationRepo := &FakeNotificationRepository{}

	service := &NotificationService{
		NotificationRepo: fakeNotificationRepo,
	}

	err := service.NotifyCommentCreated(context.Background(), notification.RecipientID, notification.ActorID, notification.CommentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fakeNotificationRepo.CreateNotificationCalled {
		t.Fatal("repository should not be called")
	}
}

func TestNotifyCommentCreated_Success(t *testing.T) {
	notification := newTestNotification()

	fakeNotificationRepo := &FakeNotificationRepository{
		CreateNotificationFunc: func(ctx context.Context, notification models.Notification) error {
			return nil
		},
	}
	service := &NotificationService{
		NotificationRepo: fakeNotificationRepo,
	}

	err := service.NotifyCommentCreated(context.Background(), notification.RecipientID, notification.ActorID, notification.CommentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !fakeNotificationRepo.CreateNotificationCalled {
		t.Fatal("repository should be called")
	}
}

func TestNotifyCommentCreated_RepositoryError(t *testing.T) {
	notification := newTestNotification()

	repoErr := errors.New("database unavailable")

	fakeNotificationRepo := &FakeNotificationRepository{
		CreateNotificationFunc: func(ctx context.Context, notification models.Notification) error {
			return repoErr
		},
	}
	service := &NotificationService{
		NotificationRepo: fakeNotificationRepo,
	}

	err := service.NotifyCommentCreated(context.Background(), notification.RecipientID, notification.ActorID, notification.CommentID)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
	if !fakeNotificationRepo.CreateNotificationCalled {
		t.Fatal("expected repository CreateNotification to be called")
	}
}

func TestGetNotificationsByUser_Success(t *testing.T) {
	expectedNotifications := []models.Notification{
		newTestNotification(),
	}

	fakeNotificationRepo := &FakeNotificationRepository{
		GetNotificationsByUserFunc: func(ctx context.Context, userID int) ([]models.Notification, error) {
			return expectedNotifications, nil
		},
	}

	service := &NotificationService{
		NotificationRepo: fakeNotificationRepo,
	}

	notifications, err := service.GetNotificationsByUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !fakeNotificationRepo.GetNotificationsByUserCalled {
		t.Fatal("expected repository GetNotificationsByUser to be called")
	}

	if !reflect.DeepEqual(notifications, expectedNotifications) {
		t.Fatalf("expected %+v, got %+v", expectedNotifications, notifications)
	}
}

func TestGetNotificationsByUser_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	fakeNotificationRepo := &FakeNotificationRepository{
		GetNotificationsByUserFunc: func(ctx context.Context, userID int) ([]models.Notification, error) {
			return nil, expectedErr
		},
	}

	service := &NotificationService{
		NotificationRepo: fakeNotificationRepo,
	}

	notifications, err := service.GetNotificationsByUser(context.Background(), 1)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected expectedErr, got %v", err)
	}

	if notifications != nil {
		t.Fatalf("expected nil notifications, got %+v", notifications)
	}

	if !fakeNotificationRepo.GetNotificationsByUserCalled {
		t.Fatal("expected repository GetNotificationsByUser to be called")
	}
}
