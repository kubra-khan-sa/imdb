# IMDb Content Upload and Review System

A content management system that allows the content team to upload movie-related data via CSV and consume it through APIs. Built with Go and MongoDB.

## Features

- **CSV Upload**: Upload movie CSV files up to 1GB (streaming parser for memory efficiency)
- **Paginated List**: View all movies with configurable pagination
- **Filtering**: Filter by release year and language
- **Sorting**: Sort by `release_date` or `vote_average` (ascending/descending)

## CSV Format

The CSV should have the following columns (comma or tab-delimited):

| Column | Description |
|--------|-------------|
| budget | Movie budget |
| homepage | Movie homepage URL |
| original_language | ISO language code (e.g. en) |
| original_title | Original title |
| overview | Movie overview/description |
| release_date | Date in YYYY-MM-DD format |
| revenue | Movie revenue |
| runtime | Runtime in minutes |
| status | Release status |
| title | Movie title |
| vote_average | Average vote/rating |
| vote_count | Number of votes |
| production_company_id | Production company ID |
| genre_id | Genre ID |
| languages | Comma-separated list of languages |

## Prerequisites

- Go 1.21+
- MongoDB (local or MongoDB Atlas)

## Setup

### 1. Clone and install dependencies

```bash
git clone <repo-url>
cd imdb
go mod download
```

### 2. Configure MongoDB

Set the MongoDB connection URI:

```bash
export MONGODB_URI="mongodb://localhost:27017"
```

For MongoDB Atlas:

```bash
export MONGODB_URI="mongodb+srv://<username>:<password>@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority"
```

**Use local MongoDB** (ignores MONGODB_URI if set):

```bash
export MONGODB_USE_LOCAL=1
# Or: unset MONGODB_URI
```

Optional: Set database name (default: `mydatabase`):

```bash
export MONGODB_DATABASE="imdb"
```

### 3. Run the server

```bash
go run cmd/server/main.go
```

The server starts on `http://localhost:8080`.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/upload` | Upload CSV file |
| GET | `/api/v1/movies` | List movies (paginated, filterable, sortable) |
| GET | `/api/v1/movies/filters` | Get available filter options (years, languages) |

### Upload CSV

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@/path/to/movies.csv"
```

**Response:**
```json
{
  "message": "File processed successfully",
  "total_processed": 5000,
  "total_inserted": 5000,
  "errors": 0
}
```

*Note: Duplicate movies (same title + release year) are updated; new movies are appended.*

### List Movies

**Query parameters:**
- `page` (default: 1) – Page number
- `per_page` (default: 10) – Items per page
- `year` – Filter by release year
- `language` – Filter by language (e.g. `en`, `ja`)
- `sort_by` – `release_date` or `vote_average`
- `sort_order` – `asc` or `desc`

**Example:**
```bash
curl "http://localhost:8080/api/v1/movies?page=1&per_page=10&year=2020&language=en&sort_by=vote_average&sort_order=desc"
```

**Response:**
```json
{
  "data": [...],
  "total": 1500,
  "page": 1,
  "per_page": 10,
  "total_pages": 150
}
```

### Get Filter Options

```bash
curl http://localhost:8080/api/v1/movies/filters
```

**Response:**
```json
{
  "years": [2020, 2021, 2022, ...],
  "languages": ["en", "ja", "es", ...]
}
```

## Postman Collection

Import the Postman collection from `postman/IMDb-Movies-API.postman_collection.json` to test all APIs:

1. Open Postman
2. File → Import → Upload the JSON file
3. The collection includes:
   - Health Check
   - Upload CSV
   - List Movies (multiple variants with filters/sort)
   - Get Filter Options

Set the `baseUrl` variable to `http://localhost:8080` (default).

## Project Structure

```
imdb/
├── cmd/server/          # Application entry point
├── internal/
│   ├── database/        # MongoDB connection
│   ├── handlers/        # HTTP handlers (upload, movies)
│   ├── models/          # Movie model
│   ├── repository/      # MongoDB repository layer
│   └── services/        # Upload service (streaming CSV)
├── pkg/csv/             # Streaming CSV parser
├── postman/             # Postman collection
└── Readme.md
```

## Running Tests

```bash
go test ./...
```
