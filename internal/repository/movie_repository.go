package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
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