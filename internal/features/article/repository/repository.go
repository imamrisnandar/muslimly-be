package repository

import (
	"muslimly-be/internal/features/article/dto"
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID            string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title         string
	Content       string
	Summary       string
	Category      string
	TitleEn       string
	ContentEn     string
	SummaryEn     string
	CategoryEn    string
	Author        string
	PublishedAt   time.Time
	ValidFrom     *time.Time
	ValidTo       *time.Time
	OrderPriority int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ArticleRepository interface {
	FindAll(filter dto.ArticleFilter) ([]Article, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) FindAll(filter dto.ArticleFilter) ([]Article, error) {
	var articles []Article
	now := time.Now()

	query := r.db.Model(&Article{}).
		Where("valid_from <= ?", now).
		Where("valid_to IS NULL OR valid_to >= ?", now)

	err := query.Order("order_priority DESC, published_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&articles).Error

	return articles, err
}
