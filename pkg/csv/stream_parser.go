package csv

import (
	"bufio"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"imdb-movies/internal/models"
)

type StreamParser struct {
	reader *csv.Reader
	now    time.Time
}

const (
	colBudget              = 0
	colHomepage            = 1
	colOriginalLanguage    = 2
	colOriginalTitle       = 3
	colOverview            = 4
	colReleaseDate         = 5
	colRevenue             = 6
	colRuntime             = 7
	colStatus              = 8
	colTitle               = 9
	colVoteAverage         = 10
	colVoteCount           = 11
	colProductionCompanyID = 12
	colGenreID             = 13
	colLanguages           = 14
)

// NewStreamParser creates a new CSV stream parser for large CSV files
func NewStreamParser(r io.Reader) *StreamParser {
	reader := csv.NewReader(bufio.NewReader(r))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // Allow variable field count
	reader.Comma = ','
	return &StreamParser{
		reader: reader,
		now:    time.Now().UTC(),
	}
}

// ReadHeader reads and returns the header row
func (p *StreamParser) ReadHeader() ([]string, error) {
	return p.reader.Read()
}

// ParseBatch reads up to batchSize rows and returns parsed movies
func (p *StreamParser) ParseBatch(batchSize int) ([]*models.Movie, error) {
	movies := make([]*models.Movie, 0, batchSize)
	for i := 0; i < batchSize; i++ {
		record, err := p.reader.Read()
		if err == io.EOF {
			return movies, nil
		}
		if err != nil {
			return nil, err
		}
		movie, err := p.ParseRow(record)
		if err != nil {
			continue // Skip invalid rows, don't fail entire batch
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

// parseDate parses release_date in YYYY-MM-DD format (e.g. 2020-12-16)
func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (p *StreamParser) ParseRow(record []string) (*models.Movie, error) {
	for len(record) < 15 {
		record = append(record, "")
	}
	releaseDate := parseDate(record[colReleaseDate])
	releaseYear := 0
	if !releaseDate.IsZero() {
		releaseYear = releaseDate.Year()
	}
	movie := &models.Movie{
		Budget:              parseInt(record[colBudget]),
		Homepage:            strings.TrimSpace(record[colHomepage]),
		OriginalLanguage:    strings.TrimSpace(record[colOriginalLanguage]),
		OriginalTitle:       strings.TrimSpace(record[colOriginalTitle]),
		Overview:            strings.TrimSpace(record[colOverview]),
		ReleaseDate:         releaseDate,
		ReleaseYear:         releaseYear,
		Revenue:             parseInt(record[colRevenue]),
		Runtime:             parseInt(record[colRuntime]),
		Status:              strings.TrimSpace(record[colStatus]),
		Title:               strings.TrimSpace(record[colTitle]),
		VoteAverage:         parseFloat(record[colVoteAverage]),
		VoteCount:           parseInt(record[colVoteCount]),
		ProductionCompanyId: strings.TrimSpace(record[colProductionCompanyID]),
		GenreId:             strings.TrimSpace(record[colGenreID]),
		Languages:           parseLanguages(record[colLanguages]),
		CreatedAt:           p.now,
		UpdatedAt:           p.now,
	}
	return movie, nil
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int(f)
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func parseLanguages(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
