package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
)

type ShowPostService struct {
	Repo *repository.ShowPostRepository
}

func (s *ShowPostService) GetPost(id string) (models.Post, error) {
	post, err := s.Repo.GetPostById(id)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}
