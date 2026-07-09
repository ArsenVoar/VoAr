package repository

import (
	"VoAr/internal/database"
	"VoAr/internal/models"
	"context"
)

type NotificationRepository struct {
	DB database.DBTX
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, notification models.Notification) error {
	_, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO notifications (recipient_id, actor_id, comment_id, created_at) VALUES ($1, $2, $3, $4)",
		notification.RecipientID, notification.ActorID, notification.CommentID, notification.CreatedAt,
	)
	return err
}

func (r *NotificationRepository) GetNotificationsByUser(ctx context.Context, userID int) ([]models.Notification, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, recipient_id, actor_id, comment_id, created_at FROM notifications WHERE recipient_id = $1 ORDER BY created_at DESC, id DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.ID, &notification.RecipientID, &notification.ActorID, &notification.CommentID, &notification.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
