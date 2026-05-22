package service

import (
	"VoAr/internal/models"
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type FakeUserRepository struct {
	User models.User
	Err  error
}

func (f *FakeUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	return f.User, f.Err
}

func (f *FakeUserRepository) GetByID(
	ctx context.Context,
	id string,
) (models.User, error) {
	return models.User{}, nil
}

func (f *FakeUserRepository) CreateUser(
	ctx context.Context,
	user models.User,
) error {
	return nil
}

func TestLogin_InvalidInput(t *testing.T) {
	fakeRepo := &FakeUserRepository{}

	service := &UserService{
		Repo: fakeRepo,
	}

	_, err := service.Login(
		context.Background(),
		"",
		"",
	)

	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf(
			"expected ErrInvalidInput, got %v",
			err,
		)
	}

}

func TestLogin_Success(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("123456"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf(
			"failed to hash password: %v",
			err,
		)
	}

	fakeRepo := &FakeUserRepository{
		User: models.User{
			Email:    "test@gmail.com",
			Password: string(hashedPassword),
		},
	}
	service := &UserService{
		Repo: fakeRepo,
	}

	user, err := service.Login(
		context.Background(),
		"test@gmail.com",
		"123456",
	)
	if err != nil {
		t.Errorf(
			"expected nil error, got %v",
			err,
		)
	}
	if user.Email != "test@gmail.com" {
		t.Errorf(
			"expected correct user email, got %v",
			user.Email,
		)
	}

}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("correct-password"),
		bcrypt.DefaultCost,
	)
	t.Fatalf(
		"failed to hash password: %v",
		err,
	)

	fakeRepo := &FakeUserRepository{
		User: models.User{
			Email:    "test@gmail.com",
			Password: string(hashedPassword),
		},
	}

	service := &UserService{
		Repo: fakeRepo,
	}

	_, err = service.Login(
		context.Background(),
		"test@gmail.com",
		"wrong-password",
	)

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf(
			"expected ErrUnauthorized, got %v",
			err,
		)
	}
}
