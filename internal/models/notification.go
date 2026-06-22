package models

import "time"

type Notification struct {
	ID          int
	RecipientID int
	ActorID     int
	CommentID   int
	CreatedAt   time.Time
}
