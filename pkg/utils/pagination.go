package utils

import (
	"strings"

	"gorm.io/gorm"
)

type PaginationConfig struct {
	Page    int
	Limit   int
	Sort    string
	Filters map[string]interface{}
}

// Paginate applies filters, sorting, and pagination to the query.
// It returns the total count of records (matching string filters) and modifies the query for pagination.
func Paginate(db *gorm.DB, cfg PaginationConfig) (*gorm.DB, int64) {
	var total int64
	query := db

	// 1. Apply Filters
	for key, value := range cfg.Filters {
		strVal, isString := value.(string)

		// 1. Use explicit operator if key contains space, >, < (e.g., "created_at >")
		if strings.Contains(key, " ") || strings.Contains(key, ">") || strings.Contains(key, "<") {
			query = query.Where(key+" ?", value)
		} else if isString && (strings.Contains(strVal, "%")) {
			// 2. Use ILIKE if value contains % (Case Insensitive)
			query = query.Where(key+" ILIKE ?", strVal)
		} else {
			// 3. Exact Match
			query = query.Where(key+" = ?", value)
		}
	}

	// 2. Count Total (before pagination/sort)
	query.Count(&total)

	// 3. Apply Sort
	if cfg.Sort != "" {
		query = query.Order(cfg.Sort)
	} else {
		// Default sort if none provided (can be overridden by caller if needed before passing DB)
		// But usually caller sets default. We'll fallback to created_at desc if strictly needed,
		// but better to leave it emptiness if not specified, or let caller handle default.
		// However, for consistency with previous code:
		query = query.Order("created_at desc")
	}

	// 4. Apply Pagination
	if cfg.Page <= 0 {
		cfg.Page = 1
	}
	if cfg.Limit <= 0 {
		cfg.Limit = 10
	}
	offset := (cfg.Page - 1) * cfg.Limit
	query = query.Offset(offset).Limit(cfg.Limit)

	return query, total
}
