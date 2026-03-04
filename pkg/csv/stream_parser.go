package csv
import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"time"
	 "imdb-movies/internal/models"
)

type NewsStreamParser struct {
	reader *csv.Reader
	now time.Time
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


// NewStreamParser creates a new CSV stream parser
func NewStreamParser(r io.Reader) *StreamParser {
    reader := csv.NewReader(bufio.NewReader(r))
    reader.LazyQuotes = true
    reader.TrimLeadingSpace = true
    reader.FieldsPerRecord = -1 // Allow variable field count
    return &StreamParser{
        reader: reader,
        now:    time.Now().UTC(),
    }
}


func (p *NewsStreamParser) ParseRow(record []string) (*models.Movie, error) {	
	for len(record) < 15 {
		record = append(record, "")
	}
	movie := &models.Movie{
		Budget:              parseFloat(record[colBudget]),
		Homepage:            record[colHomepage],
		OriginalLanguage:    record[colOriginalLanguage],
		OriginalTitle:       record[colOriginalTitle],
		Overview:              record[colOverview],
		ReleaseDate:           record[colReleaseDate],
		Revenue:               parseFloat(record[colRevenue]),
		Runtime:               parseInt(record[colRuntime]),
		Status:                record[colStatus],
		Title:                 record[colTitle],
		VoteAverage:           parseFloat(record[colVoteAverage]),
		VoteCount:             parseInt(record[colVoteCount]),
		ProductionCompanyID:   parseInt(record[colProductionCompanyID]),
		GenreID:               parseInt(record[colGenreID]),
		Languages:             record[colLanguages],
		CreatedAt:             p.now,
		UpdatedAt:             p.now,
	}
	releaseDate, err := time.Parse("2006-01-02", record[colReleaseDate])
	if err == nil {
		movie.ReleaseDate = releaseDate
	}
	movie.ReleaseYear = movie.ReleaseDate.Year()
	return movie, nil
}

