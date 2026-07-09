package service

import (
	"VoAr/internal/cache"
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"strconv"
)

type PostRepository interface {
	CreatePost(ctx context.Context, userID int, title, anons, fullText string) error
	GetPostByID(ctx context.Context, id int) (models.Post, error)
	GetPosts(ctx context.Context, page, pageSize int) ([]models.Post, error)
}
type PostService struct {
	PostRepo PostRepository
	Cache    cache.CacheService
}

const (
	defaultPage     = 1
	defaultPageSize = 10

	postCachePrefix = "post:"
)

func (s *PostService) CreatePost(ctx context.Context, userID int, title, anons, fullText string) error {
	if userID <= 0 {
		return ErrUnauthorized
	}
	if title == "" || anons == "" || fullText == "" {
		return ErrEmptyFields
	}

	return s.PostRepo.CreatePost(ctx, userID, title, anons, fullText)
}

func (s *PostService) GetPostByID(ctx context.Context, id int) (models.Post, error) {
	var post models.Post
	cacheKey := postCachePrefix + strconv.Itoa(id)

	cachedData, err := s.Cache.Get(ctx, cacheKey)
	if err == nil {
		if err := json.Unmarshal(cachedData, &post); err == nil {
			return post, nil
		}
	}

	post, err = s.PostRepo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.Post{}, ErrNotFound
		}
		return models.Post{}, err
	}

	if data, err := json.Marshal(post); err == nil {
		_ = s.Cache.Set(ctx, cacheKey, data, 0)
	}

	return post, nil
}

func (s *PostService) GetPosts(ctx context.Context, pageParam string) ([]models.Post, error) {
	page := defaultPage

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err != nil {
			return nil, err
		}

		page = p
		if page <= 0 {
			return nil, ErrInvalidInput
		}
	}

	posts, err := s.PostRepo.GetPosts(ctx, page, defaultPageSize)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
