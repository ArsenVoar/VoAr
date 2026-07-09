package service

import (
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestCreateComment_Success(t *testing.T) {
	post := newTestPost()
	comment := newTestComment()

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return post, nil
		},
	}

	fakeCommentRepo := &FakeCommentRepository{
		CreateCommentFunc: func(ctx context.Context, comment models.Comment) (int, error) {
			return 100, nil
		},
	}

	fakeNotifier := &FakeNotificationNotifier{
		NotifyFunc: func(ctx context.Context, recipientID, actorID, commentID int) error {
			return nil
		},
	}

	service := &CommentService{
		CommentRepo: fakeCommentRepo,
		PostReader:  fakePostRepo,
		Notifier:    fakeNotifier,
	}

	id, err := service.CreateComment(context.Background(), comment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 100 {
		t.Errorf("expected id 100, got %d", id)
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}
	if !fakeCommentRepo.CreateCommentCalled {
		t.Fatal("expected repository CreateComment to be called")
	}
	if !fakeNotifier.NotifyCalled {
		t.Fatal("expected notifier NotifyCommentCreated to be called")
	}
}

func TestCreateComment_EmptyContent(t *testing.T) {
	comment := newTestComment()
	comment.Content = ""

	service := &CommentService{}

	_, err := service.CreateComment(context.Background(), comment)
	if !errors.Is(err, ErrEmptyFields) {
		t.Fatalf("expected ErrEmptyFields, got %v", err)
	}
}

func TestCreateComment_PostNotFound(t *testing.T) {
	comment := newTestComment()

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return models.Post{}, repository.ErrNotFound
		},
	}
	fakeCommentRepo := &FakeCommentRepository{}

	service := &CommentService{
		PostReader:  fakePostRepo,
		CommentRepo: fakeCommentRepo,
	}

	_, err := service.CreateComment(context.Background(), comment)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}
	if fakeCommentRepo.CreateCommentCalled {
		t.Fatal("expected repository CreateComment not to be called")
	}
}

func TestCreateComment_TooLongComment(t *testing.T) {
	post := newTestPost()
	comment := newTestComment()
	comment.Content = strings.Repeat("a", maxCommentLength+1)

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return post, nil
		},
	}

	service := &CommentService{PostReader: fakePostRepo}

	_, err := service.CreateComment(context.Background(), comment)
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID not to be called")
	}
}

func TestCreateComment_RepositoryError(t *testing.T) {
	post := newTestPost()
	comment := newTestComment()

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return post, nil
		},
	}

	repoErr := errors.New("database unavailable")

	fakeCommentRepo := &FakeCommentRepository{
		CreateCommentFunc: func(ctx context.Context, comment models.Comment) (int, error) {
			return 0, repoErr
		},
	}
	fakeNotifier := &FakeNotificationNotifier{
		NotifyFunc: func(ctx context.Context, recipientID, actorID, commentID int) error {
			return nil
		},
	}

	service := &CommentService{
		CommentRepo: fakeCommentRepo,
		PostReader:  fakePostRepo,
		Notifier:    fakeNotifier,
	}

	_, err := service.CreateComment(context.Background(), comment)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}

	if !fakeCommentRepo.CreateCommentCalled {
		t.Fatal("expected repository CreateComment to be called")
	}

	if fakeNotifier.NotifyCalled {
		t.Fatal("expected notifier NotifyCommentCreated not to be called")
	}
}

func TestCreateComment_NotificationError(t *testing.T) {
	post := newTestPost()
	comment := newTestComment()

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return post, nil
		},
	}

	fakeCommentRepo := &FakeCommentRepository{
		CreateCommentFunc: func(ctx context.Context, comment models.Comment) (int, error) {
			return 100, nil
		},
	}

	notifyErr := errors.New("notification service unavailable")

	fakeNotifier := &FakeNotificationNotifier{
		NotifyFunc: func(ctx context.Context, recipientID, actorID, commentID int) error {
			return notifyErr
		},
	}

	service := &CommentService{
		CommentRepo: fakeCommentRepo,
		PostReader:  fakePostRepo,
		Notifier:    fakeNotifier,
	}

	id, err := service.CreateComment(context.Background(), comment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 100 {
		t.Errorf("expected id 100, got %d", id)
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}

	if !fakeCommentRepo.CreateCommentCalled {
		t.Fatal("expected repository CreateComment to be called")
	}

	if !fakeNotifier.NotifyCalled {
		t.Fatal("expected notifier NotifyCommentCreated to be called")
	}
}
