package handlers

import (
	"context"
	"net/http"
	"imdb-movies/internal/services"
	 "github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *services.UploadService
	maxfileSize int64
}

func NewUploadHandler(service *services.UploadService, maxFileSize int64) *UploadHandler {
	return &UploadHandler{
		service: service,
		maxfileSize: maxFileSize,
	}
	
}	

func Uploadcsv(c *gin.Context) error {
	if h.UploadService == nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload service not available"})
	}
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
	}
	if file.Size > h.maxFileSize {
		return c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the maximum allowed"})
	}
	if file.Size == 0 {
		return c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is empty"})
	}
	f, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
	}
	defer f.Close()
	movies, err := h.UploadService.ProcessUpload(c.Request.Context(), f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process uploaded file"})
	}
	return c.JSON(http.StatusOK, gin.H{
		"message":        "File processed successfully",
		"total_processed": movies.TotalProcessed,
		"total_inserted":  movies.TotalInserted,
		"errors":         movies.Errors,
	})
}	