package service

import (
	"VoAr/internal/contextkeys"
	"VoAr/internal/logger"
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
	DB   TransactionManager
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
	committed := false

	requestID := contextkeys.GetRequestID(ctx)

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	logger.TransactionStarted(requestID)

	defer func() {
		if !committed {
			tx.Rollback()
			logger.TransactionRollback(requestID, err)
		}
	}()

	txRepo := &repository.UserRepository{
		DB: tx,
	}

	_, err = txRepo.GetByEmail(ctx, user.Email)
	if err == nil {
		err = ErrUserExists
		logger.TransactionFailed(requestID, err)
		return err
	}

	if !errors.Is(err, repository.ErrNotFound) {
		logger.TransactionFailed(requestID, err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		logger.TransactionFailed(requestID, err)
		return err
	}

	user.Password = string(hashedPassword)

	err = txRepo.CreateUser(ctx, user)
	if err != nil {
		logger.TransactionFailed(requestID, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.TransactionFailed(requestID, err)
		return err
	}

	committed = true

	logger.TransactionCommitted(requestID)

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
