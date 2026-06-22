package service

import (
	"VoAr/internal/repository"
	"context"
	"database/sql"
	"log"
)

type ArticleService struct {
	Repo *repository.ArticleRepository
	DB   *sql.DB
}

func (s *ArticleService) Create(ctx context.Context, userID int, title, anons, fullText string) error {
	if userID <= 0 {
		return ErrUnauthorized
	}
	if title == "" || anons == "" || fullText == "" {
		return ErrEmptyFields
	}

	log.Printf("ArticleService userID=%d", userID)

	err := s.Repo.Create(ctx, userID, title, anons, fullText)
	if err != nil {
		return err
	}
	return nil
}
