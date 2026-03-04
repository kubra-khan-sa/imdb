package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"imdb-movies/internal/models"
	"log"
)

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) *MovieRepository {
	collection := db.Collection("movies")
	return &MovieRepository{
		collection: collection,
	}
}

func (r *MovieRepository) InsertMany(ctx context.Context, movies []*models.Movie) error {

	if len(movies) == 0 {
		return nil
	}

	docs := make([]interface{}, len(movies))
	for i, movie := range movies {
		docs[i] = movie
	}

	_, err := r.collection.InsertMany(ctx, docs)
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
		}else if y, ok := v.(int32); ok {
			years = append(years, int(y))
		}else if y, ok := v.(int64); ok {
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

func (r *MovieRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}