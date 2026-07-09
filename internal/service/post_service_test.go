package service

import (
	"VoAr/internal/cache"
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestCreatePost_Success(t *testing.T) {
	post := newTestPost()

	fakePostRepository := &FakePostRepository{
		CreatePostFunc: func(ctx context.Context, userID int, title, anons, fullText string) error {
			return nil
		},
	}

	service := PostService{
		PostRepo: fakePostRepository,
	}

	err := service.CreatePost(context.Background(), post.UserID, post.Title, post.Anons, post.FullText)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !fakePostRepository.CreatePostCalled {
		t.Fatal("expected repository CreatePost to be called")
	}
}

func TestCreatePost_Unauthorized(t *testing.T) {
	post := newTestPost()
	post.UserID = -1

	fakePostRepository := &FakePostRepository{
		CreatePostFunc: func(ctx context.Context, userID int, title, anons, fullText string) error {
			return nil
		},
	}

	service := PostService{
		PostRepo: fakePostRepository,
	}

	err := service.CreatePost(context.Background(), post.UserID, post.Title, post.Anons, post.FullText)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if fakePostRepository.CreatePostCalled {
		t.Fatal("expected repository CreatePost not to be called")
	}
}

func TestCreatePost_EmptyFields(t *testing.T) {
	post := newTestPost()
	post.Title = ""

	fakePostRepository := &FakePostRepository{
		CreatePostFunc: func(ctx context.Context, userID int, title, anons, fullText string) error {
			return nil
		},
	}

	service := PostService{
		PostRepo: fakePostRepository,
	}

	err := service.CreatePost(context.Background(), post.UserID, post.Title, post.Anons, post.FullText)
	if !errors.Is(err, ErrEmptyFields) {
		t.Fatalf("expected ErrEmptyFields, got %v", err)
	}
	if fakePostRepository.CreatePostCalled {
		t.Fatal("expected repository CreatePost not to be called")
	}
}

func TestCreatePost_RepositoryError(t *testing.T) {
	post := newTestPost()
	repoErr := errors.New("database unavailable")

	fakePostRepository := &FakePostRepository{
		CreatePostFunc: func(ctx context.Context, userID int, title, anons, fullText string) error {
			return repoErr
		},
	}

	service := PostService{
		PostRepo: fakePostRepository,
	}

	err := service.CreatePost(context.Background(), post.UserID, post.Title, post.Anons, post.FullText)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got %v", err)
	}
	if !fakePostRepository.CreatePostCalled {
		t.Fatal("expected repository CreatePost to be called")
	}
}

func TestPostService_GetPostByID_CacheHit(t *testing.T) {
	expectedPost := newTestPost()

	cachedPost, err := json.Marshal(expectedPost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return expectedPost, nil
		},
	}
	fakeCache := &FakeCacheService{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return cachedPost, nil
		},
	}

	service := &PostService{
		PostRepo: fakePostRepo,
		Cache:    fakeCache,
	}

	post, err := service.GetPostByID(context.Background(), expectedPost.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !fakeCache.GetCalled {
		t.Fatalf("expected cache Get to be called")
	}
	if fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID not to be called")
	}
	if !reflect.DeepEqual(post, expectedPost) {
		t.Fatalf("expected %+v, got %+v", expectedPost, post)
	}
}

func TestPostService_GetPostByID_CacheMiss(t *testing.T) {
	expectedPost := newTestPost()
	cacheMissErr := cache.ErrCacheMiss

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return expectedPost, nil
		},
	}
	fakeCache := &FakeCacheService{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return nil, cacheMissErr
		},
		SetFunc: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
			return nil
		},
	}

	service := &PostService{
		PostRepo: fakePostRepo,
		Cache:    fakeCache,
	}

	post, err := service.GetPostByID(context.Background(), expectedPost.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !fakeCache.GetCalled {
		t.Fatalf("expected cache Get to be called")
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}
	if !fakeCache.SetCalled {
		t.Fatalf("expected cache Set to be called")
	}
	if !reflect.DeepEqual(post, expectedPost) {
		t.Fatalf("expected %+v, got %+v", expectedPost, post)
	}
}

func TestPostService_GetPostByID_InvalidCachedData(t *testing.T) {
	expectedPost := newTestPost()

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return expectedPost, nil
		},
	}
	fakeCache := &FakeCacheService{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return []byte("{invalid json"), nil
		},
		SetFunc: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
			return nil
		},
	}

	service := &PostService{
		PostRepo: fakePostRepo,
		Cache:    fakeCache,
	}

	post, err := service.GetPostByID(context.Background(), expectedPost.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !fakeCache.GetCalled {
		t.Fatalf("expected cache Get to be called")
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}
	if !fakeCache.SetCalled {
		t.Fatalf("expected cache Set to be called")
	}
	if !reflect.DeepEqual(post, expectedPost) {
		t.Fatalf("expected %+v, got %+v", expectedPost, post)
	}
}

func TestPostService_GetPostByID_RepositoryNotFound(t *testing.T) {
	expectedPost := newTestPost()
	cacheMissErr := cache.ErrCacheMiss

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return expectedPost, repository.ErrNotFound
		},
	}
	fakeCache := &FakeCacheService{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return nil, cacheMissErr
		},
		SetFunc: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
			return nil
		},
	}

	service := &PostService{
		PostRepo: fakePostRepo,
		Cache:    fakeCache,
	}

	_, err := service.GetPostByID(context.Background(), expectedPost.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if !fakeCache.GetCalled {
		t.Fatalf("expected cache Get to be called")
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}
	if fakeCache.SetCalled {
		t.Fatalf("expected cache Set not to be called")
	}
}

func TestPostService_GetPostByID_RepositoryError(t *testing.T) {
	expectedPost := newTestPost()
	cacheMissErr := cache.ErrCacheMiss
	expectedErr := errors.New("database unavailable")

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return expectedPost, expectedErr
		},
	}
	fakeCache := &FakeCacheService{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return nil, cacheMissErr
		},
		SetFunc: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
			return nil
		},
	}

	service := &PostService{
		PostRepo: fakePostRepo,
		Cache:    fakeCache,
	}

	_, err := service.GetPostByID(context.Background(), expectedPost.ID)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected expectedErr, got %v", err)
	}
	if !fakeCache.GetCalled {
		t.Fatalf("expected cache Get to be called")
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}
	if fakeCache.SetCalled {
		t.Fatalf("expected cache Set not to be called")
	}
}

func TestPostService_GetPostByID_CacheSetError(t *testing.T) {
	expectedPost := newTestPost()
	cacheMissErr := cache.ErrCacheMiss
	cacheSetErr := cache.ErrRedisUnavailable

	fakePostRepo := &FakePostRepository{
		GetPostByIDFunc: func(ctx context.Context, id int) (models.Post, error) {
			return expectedPost, nil
		},
	}
	fakeCache := &FakeCacheService{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return nil, cacheMissErr
		},
		SetFunc: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
			return cacheSetErr
		},
	}

	service := &PostService{
		PostRepo: fakePostRepo,
		Cache:    fakeCache,
	}

	post, err := service.GetPostByID(context.Background(), expectedPost.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !fakeCache.GetCalled {
		t.Fatalf("expected cache Get to be called")
	}
	if !fakePostRepo.GetPostByIDCalled {
		t.Fatal("expected repository GetPostByID to be called")
	}
	if !fakeCache.SetCalled {
		t.Fatalf("expected cache Set to be called")
	}
	if !reflect.DeepEqual(post, expectedPost) {
		t.Fatalf("expected %+v, got %+v", expectedPost, post)
	}
}

func TestPostService_GetPosts_Success(t *testing.T) {
	pageParam := "1"
	expectedPosts := []models.Post{
		{ID: 1},
		{ID: 2},
	}

	postRepo := &FakePostRepository{
		GetPostsFunc: func(ctx context.Context, page, pageSize int) ([]models.Post, error) {
			return expectedPosts, nil
		},
	}

	service := &PostService{
		PostRepo: postRepo,
	}

	posts, err := service.GetPosts(context.Background(), pageParam)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !postRepo.GetPostsCalled {
		t.Fatal("expected repository GetPosts to be called")
	}
	if !reflect.DeepEqual(posts, expectedPosts) {
		t.Fatalf("expected %+v, got %+v", expectedPosts, posts)
	}
}

func TestPostService_GetPosts_InvalidPage(t *testing.T) {
	pageParam := "abc"
	expectedPosts := []models.Post{}

	postRepo := &FakePostRepository{
		GetPostsFunc: func(ctx context.Context, page, pageSize int) ([]models.Post, error) {
			return expectedPosts, nil
		},
	}

	service := &PostService{
		PostRepo: postRepo,
	}

	_, err := service.GetPosts(context.Background(), pageParam)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
	if postRepo.GetPostsCalled {
		t.Fatal("expected repository GetPosts not to be called")
	}
}

func TestPostService_GetPosts_RepositoryError(t *testing.T) {
	pageParam := "1"
	expectedPosts := []models.Post{}
	expectedErr := errors.New("database unavailable")

	postRepo := &FakePostRepository{
		GetPostsFunc: func(ctx context.Context, page, pageSize int) ([]models.Post, error) {
			return expectedPosts, expectedErr
		},
	}

	service := &PostService{
		PostRepo: postRepo,
	}

	posts, err := service.GetPosts(context.Background(), pageParam)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected expectedErr, got %v", err)
	}
	if posts != nil {
		t.Fatalf("expected nil posts, got %+v", posts)
	}
	if !postRepo.GetPostsCalled {
		t.Fatal("expected repository GetPosts to be called")
	}
}
