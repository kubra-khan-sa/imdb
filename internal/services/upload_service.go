package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"imdb-movies/internal/repository"
	"imdb-movies/pkg/csv"
)

const batchSize = 1000

type UploadService struct {
	repo *repository.MovieRepository
}

func NewUploadService(repo *repository.MovieRepository) *UploadService {
	return &UploadService{
		repo: repo,
	}
}

type UploadResult struct {
	TotalProcessed int `json:"total_processed"`
	TotalInserted  int `json:"total_inserted"`
	Errors        int `json:"errors"`
}

func (s *UploadService) ProcessUpload(ctx context.Context, r io.Reader, delimiter string) (*UploadResult, error) {
	parser := csv.NewStreamParser(r)
	if delimiter == "tab" || delimiter == "\t" {
		parser.SetComma('\t')
	}
	_, err := parser.ReadHeader()
	if err == io.EOF {
		return &UploadResult{
			TotalProcessed: 0,
			TotalInserted:  0,
			Errors:         0,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	result := &UploadResult{}
	var mu sync.Mutex

	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
			movies, err := parser.ParseBatch(batchSize)
			if err != nil {
				return result, fmt.Errorf("failed to parse CSV batch: %w", err)
			}
			if len(movies) == 0 {
				return result, nil
			}
			if err := s.repo.InsertMany(ctx, movies); err != nil {
				mu.Lock()
				result.Errors += len(movies)
				mu.Unlock()
				log.Printf("failed to insert batch: %v", err)
				continue
			}
			mu.Lock()
			result.TotalProcessed += len(movies)
			result.TotalInserted += len(movies)
			mu.Unlock()
		}
	}
}
