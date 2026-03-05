package repository

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"imdb-movies/internal/models"
)

type ListMoviesOptions struct {
	Page      int
	PerPage   int
	Year      *int
	Language  *string
	SortBy    string
	SortOrder string
	Years     []int
	Languages []string
}

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) *MovieRepository {
	collection := db.Collection("movies")
	repo := &MovieRepository{collection: collection}
	// Create indexes on startup (idempotent)
	if err := repo.EnsureIndexes(context.Background()); err != nil {
		log.Printf("warning: failed to ensure indexes: %v", err)
	}
	return repo
}

// EnsureIndexes creates indexes for efficient queries and upserts
func (r *MovieRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "title", Value: 1}, {Key: "release_year", Value: 1}},
			Options: options.Index().SetName("title_release_year"),
		},
		{Keys: bson.D{{Key: "release_year", Value: 1}}, Options: options.Index().SetName("release_year")},
		{Keys: bson.D{{Key: "original_language", Value: 1}}, Options: options.Index().SetName("original_language")},
		{Keys: bson.D{{Key: "languages", Value: 1}}, Options: options.Index().SetName("languages")},
		{Keys: bson.D{{Key: "release_date", Value: 1}}, Options: options.Index().SetName("release_date")},
		{Keys: bson.D{{Key: "vote_average", Value: 1}}, Options: options.Index().SetName("vote_average")},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// UpsertMany inserts or updates movies to avoid duplicates. Uses (title, release_year) as the unique key.
func (r *MovieRepository) UpsertMany(ctx context.Context, movies []*models.Movie) error {
	if len(movies) == 0 {
		return nil
	}
	models := make([]mongo.WriteModel, 0, len(movies))
	for _, movie := range movies {
		filter := bson.M{
			"title":        movie.Title,
			"release_year": movie.ReleaseYear,
		}
		models = append(models, mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(movie).
			SetUpsert(true))
	}
	_, err := r.collection.BulkWrite(ctx, models)
	return err
}
func (r *MovieRepository) GetDistinctYears(ctx context.Context) ([]int, error) {
	values, err := r.collection.Distinct(ctx, "release_year", bson.M{})
	if err != nil {
		return nil, err
	}
	years := make([]int, 0, len(values))
	for i, v := range values {
		if y, ok := v.(int); ok {
			years = append(years, y)
		} else if y, ok := v.(int32); ok {
			years = append(years, int(y))
		} else if y, ok := v.(int64); ok {
			years = append(years, int(y))
		} else {
			log.Printf("unexpected type for release_year at index %d: %T", i, v)
		}
	}
	return years, nil
}
func (r *MovieRepository) DeleteAll(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{})
	return err
}

func (r *MovieRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

func (r *MovieRepository) ListMovies(ctx context.Context, opts *ListMoviesOptions) (*models.MovieResponse, error) {
	filter := bson.M{}
	if opts.Year != nil && *opts.Year > 0 {
		filter["release_year"] = *opts.Year
	}
	if opts.Language != nil && *opts.Language != "" {
		filter["$or"] = []bson.M{
			{"original_language": *opts.Language},
			{"languages": *opts.Language},
		}
	}

	// Sort
	sortField := "release_date"
	if opts.SortBy == "vote_average" {
		sortField = "vote_average"
	}
	sortVal := 1
	if opts.SortOrder == "desc" {
		sortVal = -1
	}
	sortOpt := options.Find().SetSort(bson.D{{Key: sortField, Value: sortVal}})

	// Pagination
	skip := (opts.Page - 1) * opts.PerPage
	if skip < 0 {
		skip = 0
	}
	sortOpt.SetSkip(int64(skip)).SetLimit(int64(opts.PerPage))

	cur, err := r.collection.Find(ctx, filter, sortOpt)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var movies []*models.Movie
	if err := cur.All(ctx, &movies); err != nil {
		return nil, err
	}

	total, err := r.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / opts.PerPage
	if int(total)%opts.PerPage > 0 {
		totalPages++
	}

	return &models.MovieResponse{
		Data:       movies,
		Total:      int(total),
		Page:       opts.Page,
		PerPage:    opts.PerPage,
		TotalPages: totalPages,
	}, nil
}

// isValidLanguage filters out garbage values (brackets, stray quotes, etc.)
func isValidLanguage(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || len(s) < 2 {
		return false
	}
	// Reject values with brackets or Python list artifacts
	if strings.Contains(s, "[") || strings.Contains(s, "]") ||
		strings.HasPrefix(s, "'") || strings.HasSuffix(s, "'") {
		return false
	}
	// Reject placeholder/unknown
	if s == "?????" || s == "??????" || s == "''" {
		return false
	}
	return true
}

func (r *MovieRepository) GetDistinctLanguages(ctx context.Context) ([]string, error) {
	values, err := r.collection.Distinct(ctx, "languages", bson.M{})
	if err != nil {
		return nil, err
	}
	origValues, _ := r.collection.Distinct(ctx, "original_language", bson.M{})
	seen := make(map[string]bool)
	languages := make([]string, 0)
	for _, v := range values {
		if s, ok := v.(string); ok && isValidLanguage(s) && !seen[s] {
			seen[s] = true
			languages = append(languages, s)
		}
	}
	for _, v := range origValues {
		if s, ok := v.(string); ok && isValidLanguage(s) && !seen[s] {
			seen[s] = true
			languages = append(languages, s)
		}
	}
	return languages, nil
}
