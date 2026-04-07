package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"strconv"
)

type PostService struct {
	Repo *repository.PostRepository
}

func (s *PostService) GetPosts(pageParam string) ([]models.Post, error) {
	page := 1
	pageSize := 10

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err != nil {
			return nil, err
		}
		page = p
	}

	posts, err := s.Repo.GetPosts(page, pageSize)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
