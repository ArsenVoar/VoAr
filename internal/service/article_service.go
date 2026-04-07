package service

import (
	"VoAr/internal/repository"
	"errors"
)

type ArticleService struct {
	Repo *repository.ArticleRepository
}

func (s *ArticleService) Create(title, anons, fullText string) error {
	if title == "" || anons == "" || fullText == "" {
		return errors.New("empty fields")
	}

	err := s.Repo.Create(title, anons, fullText)
	if err != nil {
		return err
	}
	return nil
}
