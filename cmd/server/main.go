package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Initialize database connection
	db, err := ConnectMongoDB(ctx, "mydatabase")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")
	movierepo := NewMovieRepository(db)
	movieHandler := handlers.NewMovieHandler(movierepo)
	router :=setupRoutes(movieHandler, 10)	
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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

func setupRoutes() {
	router:=gin.Default()
	router.MaxMultipartMemory = maxUploadSize
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	v1 := router.Group("/api/v1")
	{
		v1.POST("/upload", uploadHandler)
		v1.GET("/movies", getFileHandler)
		v1.GET("/movies/filters", getFiltersHandler)
	}
}