package service

import (
	"muslimly-be/internal/features/article/dto"
	"muslimly-be/internal/features/article/repository"
)

type ArticleService interface {
	GetArticles(limit, offset int) ([]dto.ArticleResponse, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) GetArticles(limit, offset int) ([]dto.ArticleResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	filter := dto.ArticleFilter{
		Limit:  limit,
		Offset: offset,
	}

	articles, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	var responses []dto.ArticleResponse
	for _, article := range articles {
		responses = append(responses, dto.ArticleResponse{
			ID:            article.ID,
			Title:         article.Title,
			Content:       article.Content,
			Summary:       article.Summary,
			Category:      article.Category,
			Author:        article.Author,
			PublishedAt:   article.PublishedAt,
			OrderPriority: article.OrderPriority,
		})
	}

	return responses, nil
}
