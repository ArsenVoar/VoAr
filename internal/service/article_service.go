package service

import (
	"VoAr/internal/repository"
	"context"
	"database/sql"
)

type ArticleService struct {
	Repo *repository.ArticleRepository
	DB   *sql.DB
}

func (s *ArticleService) Create(ctx context.Context, title, anons, fullText string) error {
	if title == "" || anons == "" || fullText == "" {
		return ErrEmptyFields
	}

	err := s.Repo.Create(ctx, title, anons, fullText)
	if err != nil {
		return err
	}
	return nil
}
