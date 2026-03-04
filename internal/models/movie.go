package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Budget   int                `bson:"budget"`
	Homepage string             `bson:"homepage"`
	OriginalTitle string             `bson:"original_title"`
	OriginalLanguage string          `bson:"original_language"`
	Overview string             `bson:"overview"`
	ReleaseDate time.Time          `bson:"release_date"`
	ReleaseYear int                `bson:"release_year"`
	Revenue  int                `bson:"revenue"`
	Runtime  int                `bson:"runtime"`
	Status   string             `bson:"status"`
	Title	string             `bson:"title"`
	VoteAverage float64          `bson:"vote_average"`
	VoteCount int                `bson:"vote_count"`
	ProducionCompanyId string             `bson:"production_company_id"`
	GenreId string             `bson:"genre_id"`
	Languages  []string            `bson:"languages"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type MovieResponse struct {
	Data []*Movie `json:"data"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	PerPage int      `json:"per_page"`
	TotalPages int      `json:"total_pages"`
}