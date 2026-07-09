package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestUserService_Login_Success(t *testing.T) {
	user := newTestUser("123456")

	fakeUserRepo := &FakeUserRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.User, error) {
			return user, nil
		},
	}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	login, err := service.Login(context.Background(), "test@gmail.com", "123456")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if login.Email != "test@gmail.com" {
		t.Errorf("expected correct user email, got %v", user.Email)
	}
	if !fakeUserRepo.GetUserByEmailCalled {
		t.Fatal("expected repository GetUserByEmail to be called")
	}
}

func TestUserService_Login_InvalidInput(t *testing.T) {
	fakeUserRepo := &FakeUserRepository{}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	_, err := service.Login(
		context.Background(),
		"",
		"",
	)

	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
	if fakeUserRepo.GetUserByEmailCalled {
		t.Fatal("expected repository GetUserByEmail not to be called")
	}
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	user := newTestUser("123456")

	fakeUserRepo := &FakeUserRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.User, error) {
			return user, nil
		},
	}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	_, err := service.Login(
		context.Background(),
		"test@gmail.com",
		"wrong-password",
	)

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
	if !fakeUserRepo.GetUserByEmailCalled {
		t.Fatal("expected repository GetUserByEmail to be called")
	}
}

func TestUserService_Login_RepositoryError(t *testing.T) {
	user := newTestUser("123456")

	fakeUserRepo := &FakeUserRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.User, error) {
			return user, repository.ErrNotFound
		},
	}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	_, err := service.Login(
		context.Background(),
		"test@gmail.com",
		"123456",
	)

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if !fakeUserRepo.GetUserByEmailCalled {
		t.Fatal("expected repository GetUserByEmail to be called")
	}
}

func TestUserService_GetProfile_Success(t *testing.T) {
	user := newTestUser("")

	fakeUserRepo := &FakeUserRepository{
		GetUserByIDFunc: func(tx context.Context, id int) (models.User, error) {
			return user, nil
		},
	}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	profile, err := service.GetUserProfile(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(profile, user) {
		t.Fatalf("expected %+v, got %+v", user, profile)
	}
	if !fakeUserRepo.GetUserByIDCalled {
		t.Fatal("expected repository GetUserByID to be called")
	}
}

func TestUserService_GetProfile_UserNotFound(t *testing.T) {
	user := newTestUser("")

	fakeUserRepo := &FakeUserRepository{
		GetUserByIDFunc: func(tx context.Context, id int) (models.User, error) {
			return user, repository.ErrNotFound
		},
	}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	_, err := service.GetUserProfile(context.Background(), 1, 1)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if !fakeUserRepo.GetUserByIDCalled {
		t.Fatal("expected repository GetUserByID to be called")
	}
}

func TestUserService_GetProfile_RepositoryError(t *testing.T) {
	user := newTestUser("")

	expectedErr := errors.New("database unavailable")

	fakeUserRepo := &FakeUserRepository{
		GetUserByIDFunc: func(tx context.Context, id int) (models.User, error) {
			return user, expectedErr
		},
	}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	_, err := service.GetUserProfile(context.Background(), 1, 1)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected expectedErr, got %v", err)
	}
	if !fakeUserRepo.GetUserByIDCalled {
		t.Fatal("expected repository GetUserByID to be called")
	}
}

func TestUserService_GetProfile_UnauthorizedSession(t *testing.T) {
	user := newTestUser("")

	fakeUserRepo := &FakeUserRepository{
		GetUserByIDFunc: func(ctx context.Context, id int) (models.User, error) {
			return user, nil
		},
	}

	service := &UserService{
		UserRepo: fakeUserRepo,
	}

	_, err := service.GetUserProfile(context.Background(), 1, 0)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if fakeUserRepo.GetUserByIDCalled {
		t.Fatal("expected repository GetUserByID not to be called")
	}
}
