package service

import (
	"muslimly-be/internal/features/article/dto"
	"muslimly-be/internal/features/article/repository"
)

type ArticleService interface {
	GetArticles(limit, offset int, lang string) ([]dto.ArticleResponse, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) GetArticles(limit, offset int, lang string) ([]dto.ArticleResponse, error) {
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
		// Default to Indonesian
		title := article.Title
		content := article.Content
		summary := article.Summary
		category := article.Category

		// Override with English if requested and available
		if lang == "en" {
			if article.TitleEn != "" {
				title = article.TitleEn
			}
			if article.ContentEn != "" {
				content = article.ContentEn
			}
			if article.SummaryEn != "" {
				summary = article.SummaryEn
			}
			if article.CategoryEn != "" {
				category = article.CategoryEn
			}
		}

		responses = append(responses, dto.ArticleResponse{
			ID:            article.ID,
			Title:         title,
			Content:       content,
			Summary:       summary,
			Category:      category,
			Author:        article.Author,
			PublishedAt:   article.PublishedAt,
			OrderPriority: article.OrderPriority,
		})
	}

	return responses, nil
}
