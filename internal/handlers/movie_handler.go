package handlers

import (
	"context"
	"net/http"
	"imdb-movies/internal/repository"
)

type MovieHandler struct {
	repo *repository.MovieRepository
}

func NewMovieHandler(repo *repository.MovieRepository) *MovieHandler {
	return &MovieHandler{
		repo: repo,
	}
}

func (h *MovieHandler) ListMovies(c *gin.Context) error {
	opts := &repository.ListMoviesOptions{
		Page:    getIntparam(c, "page", 1),
		PerPage: getIntparam(c, "per_page", 10),
		sortBy: c.DefaultQuery("sort_by", "release_date"),
		sortOrder: c.DefaultQuery("sort_order", "desc"),
	}
	if year := c.Query("year"); year != "" {
		year, err := strconv.Atoi(year)
		if err ==nil & year > 0 {
			opts.Year = &year
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
		return c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list movies"})
	}
	return c.JSON(http.StatusOK, response)
}
