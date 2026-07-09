package service

import (
	"VoAr/internal/models"
	"context"
	"database/sql"
	"time"
)

type FakeCommentRepository struct {
	CreateCommentFunc       func(ctx context.Context, comment models.Comment) (int, error)
	GetCommentsByPostFunc   func(ctx context.Context, postID int) ([]models.Comment, error)
	CreateCommentCalled     bool
	GetCommentsByPostCalled bool
}
type FakePostRepository struct {
	CreatePostFunc    func(ctx context.Context, userID int, title, anons, fullText string) error
	GetPostByIDFunc   func(ctx context.Context, id int) (models.Post, error)
	GetPostsFunc      func(ctx context.Context, page, pageSize int) ([]models.Post, error)
	CreatePostCalled  bool
	GetPostByIDCalled bool
	GetPostsCalled    bool
}
type FakeUserRepository struct {
	GetUserByEmailFunc func(ctx context.Context, email string) (models.User, error)
	GetUserByIDFunc    func(ctx context.Context, id int) (models.User, error)
	CreateUserFunc     func(ctx context.Context, user models.User) error

	GetUserByEmailCalled bool
	GetUserByIDCalled    bool
	CreateUserCalled     bool
}
type FakeTransactionManager struct {
	Err error
}
type FakeNotificationRepository struct {
	CreateNotificationFunc       func(ctx context.Context, notification models.Notification) error
	GetNotificationsByUserFunc   func(ctx context.Context, userID int) ([]models.Notification, error)
	CreateNotificationCalled     bool
	GetNotificationsByUserCalled bool
}
type FakeNotificationNotifier struct {
	NotifyFunc   func(ctx context.Context, recipientID int, actorID int, commentID int) error
	NotifyCalled bool
}
type FakeCacheService struct {
	GetFunc      func(ctx context.Context, key string) ([]byte, error)
	SetFunc      func(ctx context.Context, key string, value []byte, ttl time.Duration) error
	DeleteFunc   func(ctx context.Context, key string) error
	TTLFunc      func(ctx context.Context, key string) (time.Duration, error)
	GetCalled    bool
	SetCalled    bool
	DeleteCalled bool
	TTLCalled    bool
}

func (f *FakeCommentRepository) CreateComment(ctx context.Context, comment models.Comment) (int, error) {
	f.CreateCommentCalled = true

	if f.CreateCommentFunc != nil {
		return f.CreateCommentFunc(ctx, comment)
	}

	panic("FakeCommentRepository.CreateCommentFunc is nil")
}

func (f *FakeCommentRepository) GetCommentsByPost(ctx context.Context, postID int) ([]models.Comment, error) {
	f.GetCommentsByPostCalled = true

	if f.GetCommentsByPostFunc != nil {
		return f.GetCommentsByPostFunc(ctx, postID)
	}

	panic("FakeCommentRepository.GetCommentsByPostFunc is nil")
}

func (f *FakeNotificationNotifier) NotifyCommentCreated(ctx context.Context, recipientID int, actorID int, commentID int) error {
	f.NotifyCalled = true

	if f.NotifyFunc != nil {
		return f.NotifyFunc(ctx, recipientID, actorID, commentID)
	}

	panic("FakeNotificationNotifier.NotifyFunc is nil")
}

func (f *FakeNotificationRepository) CreateNotification(ctx context.Context, notification models.Notification) error {
	f.CreateNotificationCalled = true

	if f.CreateNotificationFunc != nil {
		return f.CreateNotificationFunc(ctx, notification)
	}

	panic("FakeNotificationRepository.CreateNotificationFunc is nil")
}

func (f *FakeNotificationRepository) GetNotificationsByUser(ctx context.Context, userID int) ([]models.Notification, error) {
	f.GetNotificationsByUserCalled = true

	if f.GetNotificationsByUserFunc != nil {
		return f.GetNotificationsByUserFunc(ctx, userID)
	}

	panic("FakeNotificationRepository.GetNotificationsByUserFunc is nil")
}

func (f *FakePostRepository) CreatePost(ctx context.Context, userID int, title, anons, fullText string) error {
	f.CreatePostCalled = true

	if f.CreatePostFunc != nil {
		return f.CreatePostFunc(ctx, userID, title, anons, fullText)
	}

	panic("FakePostRepository.CreatePostFunc is nil")
}

func (f *FakePostRepository) GetPostByID(ctx context.Context, id int) (models.Post, error) {
	f.GetPostByIDCalled = true

	if f.GetPostByIDFunc != nil {
		return f.GetPostByIDFunc(ctx, id)
	}

	panic("FakePostRepository.GetPostByIDFunc is nil")
}

func (f *FakePostRepository) GetPosts(ctx context.Context, page, pageSize int) ([]models.Post, error) {
	f.GetPostsCalled = true

	if f.GetPostsFunc != nil {
		return f.GetPostsFunc(ctx, page, pageSize)
	}

	panic("FakePostRepository.GetPostsFunc is nil")
}

func (f *FakeCacheService) Get(ctx context.Context, key string) ([]byte, error) {
	f.GetCalled = true

	if f.GetFunc != nil {
		return f.GetFunc(ctx, key)
	}

	panic("FakeCacheService.GetFunc is nil")
}

func (f *FakeCacheService) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	f.SetCalled = true

	if f.SetFunc != nil {
		return f.SetFunc(ctx, key, value, ttl)
	}

	panic("FakeCacheService.SetFunc is nil")
}

func (f *FakeCacheService) Delete(ctx context.Context, key string) error {
	f.DeleteCalled = true

	if f.DeleteFunc != nil {
		return f.DeleteFunc(ctx, key)
	}

	panic("FakeCacheService.DeleteFunc is nil")
}

func (f *FakeCacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	f.TTLCalled = true

	if f.TTLFunc != nil {
		return f.TTLFunc(ctx, key)
	}

	panic("FakeCacheService.TTLFunc is nil")
}

func (f *FakeUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	f.GetUserByEmailCalled = true

	if f.GetUserByEmailFunc != nil {
		return f.GetUserByEmailFunc(ctx, email)
	}

	panic("FakeUserRepository.GetUserByEmailFunc is nil")
}

func (f *FakeUserRepository) GetUserByID(ctx context.Context, id int) (models.User, error) {
	f.GetUserByIDCalled = true

	if f.GetUserByIDFunc != nil {
		return f.GetUserByIDFunc(ctx, id)
	}

	panic("FakeUserRepository.GetUserByIDFunc is nil")
}

func (f *FakeUserRepository) CreateUser(ctx context.Context, user models.User) error {
	f.CreateUserCalled = true

	if f.CreateUserFunc != nil {
		return f.CreateUserFunc(ctx, user)
	}

	panic("FakeUserRepository.CreateUserFunc is nil")
}

func (f *FakeTransactionManager) BeginTx(ctx context.Context, pts *sql.TxOptions) (*sql.Tx, error) {
	return nil, f.Err
}
