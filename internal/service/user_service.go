package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"errors"
)

var ErrUserExists = errors.New("user already exists")

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) GetProfile(userId string, sessionUserID string) (models.User, error) {
	user, err := s.Repo.GetByID(userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	if sessionUserID == "" {
		return models.User{}, ErrUnauthorized
	}

	return user, nil
}

func (s *UserService) Create(user models.User) error {
	err := s.Repo.Create(user)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return ErrUserExists
		}
		return err
	}
	return nil
}
