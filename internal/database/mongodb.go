package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
// ConnectMongoDB establishes a connection to MongoDB and returns a database instance.
func ConnectMongoDB(ctx context.Context, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
   clientOptions := options.Client().ApplyURI("mongodb+srv://saba91116_db_user:<db_password>@cluster0.n0zvaw4.mongodb.net/?appName=Cluster0")

	// clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(database), nil
}