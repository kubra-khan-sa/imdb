package handlers

import (
	"net/http"

	"imdb-movies/internal/services"
	"github.com/gin-gonic/gin"
)

const maxUploadSize = 1024 * 1024 * 1024 // 1GB

type UploadHandler struct {
	service     *services.UploadService
	maxFileSize int64
}

func NewUploadHandler(service *services.UploadService, maxFileSize int64) *UploadHandler {
	return &UploadHandler{
		service:     service,
		maxFileSize: maxFileSize,
	}
}

func (h *UploadHandler) UploadCSV(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload service not available"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file: " + err.Error()})
		return
	}
	if file.Size > h.maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the maximum allowed (1GB)"})
		return
	}
	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is empty"})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer f.Close()

	result, err := h.service.ProcessUpload(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process uploaded file: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":          "File processed successfully",
		"total_processed":  result.TotalProcessed,
		"total_inserted":   result.TotalInserted,
		"errors":           result.Errors,
	})
}
