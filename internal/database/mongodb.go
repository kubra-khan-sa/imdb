package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongoDB establishes a connection to MongoDB and returns a database instance.
func ConnectMongoDB(ctx context.Context, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	if os.Getenv("MONGODB_USE_LOCAL") == "1" || uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// For Atlas: add longer timeouts to avoid server selection timeout during bulk ops
	if strings.Contains(uri, "mongodb+srv://") {
		sep := "?"
		if strings.Contains(uri, "?") {
			sep = "&"
		}
		uri = uri + sep + "serverSelectionTimeoutMS=60000&connectTimeoutMS=10000"
	}

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(database), nil
}
