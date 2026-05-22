package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) error
}

type UserService struct {
	Repo UserRepository
}

func (s *UserService) GetProfile(ctx context.Context, userId string, sessionUserID string) (models.User, error) {
	user, err := s.Repo.GetByID(ctx, userId)
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

func (s *UserService) Register(ctx context.Context, user models.User) error {
	if user.Email == "" || user.Password == "" {
		return ErrInvalidInput
	}

	_, err := s.Repo.GetByEmail(ctx, user.Email)
	if err == nil {
		return ErrUserExists
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	err = s.Repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (models.User, error) {
	if email == "" || password == "" {
		return models.User{}, ErrInvalidInput
	}

	user, err := s.Repo.GetByEmail(ctx, email)
	if err != nil {
		return models.User{}, ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return models.User{}, ErrUnauthorized
	}

	return user, nil
}
