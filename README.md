# Go Movie API

A Go API that integrates with the OMDb API to provide movie information, TV episode details, genre-based movie searches, and movie recommendations.

## Requirements

- Go 1.25.1 or higher
- OMDb API key (get one free at http://www.omdbapi.com/apikey.aspx)

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd postman_backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create a `.env` file in the root directory:
```
OMDB_API_KEY=your_omdb_api_key_here
PORT=8080
```

4. Run the application:
```bash
go run main.go
```

The server will start on port 8080 by default.

## API Endpoints

### Movie Details
```
GET /api/movie?title=<movie_title>
```
Example: `http://localhost:8080/api/movie?title=The Matrix`

### TV Episode Details
```
GET /api/episode?series_title=<series>&season=<number>&episode_number=<number>
```
Example: `http://localhost:8080/api/episode?series_title=Breaking Bad&season=1&episode_number=1`

### Movies by Genre
```
GET /api/movies/genre?genre=<genre>
```
Example: `http://localhost:8080/api/movies/genre?genre=Action`

### Movie Recommendations
```
GET /api/recommendations?favorite_movie=<movie_title>
```
Example: `http://localhost:8080/api/recommendations?favorite_movie=The Matrix`

### Health Check
```
GET /health
```

## Building

To build the application:
```bash
go build -o bin/movie-api main.go
```

To run the built binary:
```bash
./bin/movie-api
```

## Project Structure

```
go-api/
├── main.go              # Application entry point
├── go.mod              # Go module definition
├── .env                # Environment variables
├── .gitignore          # Git ignore file
├── models/
│   └── movie.go        # Data models
├── services/
│   └── omdb.go         # OMDb API service
└── handlers/
    └── movie.go        # HTTP handlers
```

## Environment Variables

- `OMDB_API_KEY`: Your OMDb API key (required)
- `PORT`: Server port (optional, defaults to 8080)
