package service

import (
	"VoAr/internal/repository"
	"context"
)

type ArticleService struct {
	Repo *repository.ArticleRepository
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
