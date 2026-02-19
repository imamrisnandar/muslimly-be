package handler

import (
	"net/http"
	"strconv"

	"muslimly-be/internal/features/article/service"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

// GetArticles godoc
// @Summary      Get Articles
// @Description  Get list of articles with pagination, valid date filter, and priority sorting
// @Tags         Article
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Limit (default 10)"
// @Param        offset  query     int  false  "Offset (default 0)"
// @Param        lang    query     string  false  "Language (default 'id')"
// @Param        search  query     string  false  "Search keyword"
// @Success      200     {array}   dto.ArticleResponse
// @Router       /articles [get]
func (h *ArticleHandler) GetArticles(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	lang := c.QueryParam("lang")
	search := c.QueryParam("search")

	articles, err := h.service.GetArticles(limit, offset, lang, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch articles"})
	}

	return c.JSON(http.StatusOK, articles)
}
