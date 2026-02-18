package dto

import "time"

type ArticleResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Summary       string    `json:"summary"`
	Category      string    `json:"category"`
	Author        string    `json:"author"`
	PublishedAt   time.Time `json:"published_at"`
	OrderPriority int       `json:"order_priority"`
}

type ArticleFilter struct {
	Limit  int
	Offset int
}
