package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"imdb-movies/internal/database"
	"imdb-movies/internal/handlers"
	"imdb-movies/internal/repository"
	"imdb-movies/internal/services"
)

const maxUploadSize = 1024 * 1024 * 1024 // 1GB

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection
	dbName := getEnv("MONGODB_DATABASE", "mydatabase")
	db, err := database.ConnectMongoDB(ctx, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	// Initialize repositories and services
	movieRepo := repository.NewMovieRepository(db)
	uploadService := services.NewUploadService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieRepo)
	uploadHandler := handlers.NewUploadHandler(uploadService, maxUploadSize)

	router := setupRoutes(movieHandler, uploadHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Minute,  // Large CSV uploads can take a long time
		WriteTimeout: 30 * time.Minute,  // Processing + DB writes for 1GB file
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Listen for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func setupRoutes(movieHandler *handlers.MovieHandler, uploadHandler *handlers.UploadHandler) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = maxUploadSize

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/upload", uploadHandler.UploadCSV)
		v1.GET("/movies", movieHandler.ListMovies)
		v1.DELETE("/movies", movieHandler.DeleteAll)
		v1.GET("/movies/filters", movieHandler.GetFilterOptions)
	}

	return router
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
