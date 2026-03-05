package handlers

import (
	"context"
	"net/http"
	"strconv"

	"imdb-movies/internal/repository"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	repo *repository.MovieRepository
}

func NewMovieHandler(repo *repository.MovieRepository) *MovieHandler {
	return &MovieHandler{
		repo: repo,
	}
}

func (h *MovieHandler) ListMovies(c *gin.Context) {
	opts := &repository.ListMoviesOptions{
		Page:      getIntParam(c, "page", 1),
		PerPage:   getIntParam(c, "per_page", 10),
		SortBy:    c.DefaultQuery("sort_by", "release_date"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}
	if year := c.Query("year"); year != "" {
		y, err := strconv.Atoi(year)
		if err == nil && y > 0 {
			opts.Year = &y
		}
	}
	if language := c.Query("language"); language != "" {
		opts.Language = &language
	}
	if opts.SortBy != "release_date" && opts.SortBy != "vote_average" {
		opts.SortBy = "release_date"
	}
	if opts.SortOrder != "asc" && opts.SortOrder != "desc" {
		opts.SortOrder = "desc"
	}

	response, err := h.repo.ListMovies(context.Background(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list movies: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *MovieHandler) DeleteAll(c *gin.Context) {
	ctx := c.Request.Context()
	if err := h.repo.DeleteAll(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movies: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All movies deleted successfully"})
}

func (h *MovieHandler) GetFilterOptions(c *gin.Context) {
	ctx := c.Request.Context()
	years, err := h.repo.GetDistinctYears(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get years"})
		return
	}
	languages, err := h.repo.GetDistinctLanguages(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get languages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"years":     years,
		"languages": languages,
	})
}

func getIntParam(c *gin.Context, name string, defaultValue int) int {
	valueStr := c.DefaultQuery(name, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
