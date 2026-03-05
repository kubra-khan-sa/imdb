package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Budget               int                `bson:"budget" json:"budget"`
	Homepage             string             `bson:"homepage" json:"homepage"`
	OriginalTitle        string             `bson:"original_title" json:"original_title"`
	OriginalLanguage     string             `bson:"original_language" json:"original_language"`
	Overview             string             `bson:"overview" json:"overview"`
	ReleaseDate          time.Time          `bson:"release_date" json:"release_date"`
	ReleaseYear          int                `bson:"release_year" json:"release_year"`
	Revenue              int                `bson:"revenue" json:"revenue"`
	Runtime              int                `bson:"runtime" json:"runtime"`
	Status               string             `bson:"status" json:"status"`
	Title                string             `bson:"title" json:"title"`
	VoteAverage          float64            `bson:"vote_average" json:"vote_average"`
	VoteCount            int                `bson:"vote_count" json:"vote_count"`
	ProductionCompanyId  string             `bson:"production_company_id" json:"production_company_id"`
	GenreId              string             `bson:"genre_id" json:"genre_id"`
	Languages            []string           `bson:"languages" json:"languages"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
}

type MovieResponse struct {
	Data []*Movie `json:"data"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	PerPage int      `json:"per_page"`
	TotalPages int      `json:"total_pages"`
}