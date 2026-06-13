package service

import (
	"VoAr/internal/cache"
	"VoAr/internal/models"
	"VoAr/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
)

type PostService struct {
	Repo  *repository.PostRepository
	DB    *sql.DB
	Cache cache.CacheService
}

func (s *PostService) GetPost(ctx context.Context, id int) (models.Post, error) {
	var post models.Post
	cacheKey := "post:" + strconv.Itoa(id)

	cachedData, err := s.Cache.Get(ctx, cacheKey)
	if err == nil {
		if err := json.Unmarshal(cachedData, &post); err == nil {
			return post, nil
		}
	}

	post, err = s.Repo.GetPostById(ctx, id)
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
	page := 1
	pageSize := 10

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err != nil {
			return nil, err
		}
		page = p
	}

	posts, err := s.Repo.GetPosts(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
