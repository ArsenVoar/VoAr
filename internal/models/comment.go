package models

import "time"

type Comment struct {
	ID        int
	ArticleID int
	UserID    int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
