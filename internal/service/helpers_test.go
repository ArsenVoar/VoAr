package service

import (
	"VoAr/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func newTestUser(password string) models.User {
	user := models.User{
		ID:    1,
		Name:  "Test",
		Email: "test@gmail.com",
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			panic(err)
		}

		user.Password = string(hashedPassword)
	}

	return user
}

func newTestPost() models.Post {
	return models.Post{
		ID:       1,
		UserID:   1,
		Title:    "Test title",
		Anons:    "Test anons",
		FullText: "Test full text",
	}
}

func newTestComment() models.Comment {
	return models.Comment{
		ID:      1,
		Content: "hello world",
		PostID:  1,
		UserID:  2,
	}
}

func newTestNotification() models.Notification {
	return models.Notification{
		ID:          1,
		RecipientID: 1,
		ActorID:     2,
		CommentID:   3,
	}
}
