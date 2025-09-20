package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"go-api/models"
)

const OMDbBaseURL = "http://www.omdbapi.com/"

type OMDbService struct {
	APIKey string
	Client *http.Client
}

func NewOMDbService(apiKey string) *OMDbService {
	return &OMDbService{
		APIKey: apiKey,
		Client: &http.Client{},
	}
}

// GetMovieByTitle fetches movie details by title
func (s *OMDbService) GetMovieByTitle(title string) (*models.OMDbResponse, error) {
	params := url.Values{}
	params.Add("apikey", s.APIKey)
	params.Add("t", title)
	params.Add("plot", "full")

	return s.makeRequest(params)
}

// GetEpisodeDetails fetches TV episode details
func (s *OMDbService) GetEpisodeDetails(seriesTitle string, season, episode int) (*models.OMDbResponse, error) {
	params := url.Values{}
	params.Add("apikey", s.APIKey)
	params.Add("t", seriesTitle)
	params.Add("Season", strconv.Itoa(season))
	params.Add("Episode", strconv.Itoa(episode))

	return s.makeRequest(params)
}

// SearchMovies searches for movies by title
func (s *OMDbService) SearchMovies(query string, page int) (*models.SearchResponse, error) {
	params := url.Values{}
	params.Add("apikey", s.APIKey)
	params.Add("s", query)
	params.Add("type", "movie")
	if page > 0 {
		params.Add("page", strconv.Itoa(page))
	}

	resp, err := http.Get(OMDbBaseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var searchResp models.SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &searchResp, nil
}

// GetMoviesByGenre collects movies of a specific genre
func (s *OMDbService) GetMoviesByGenre(genre string, limit int) ([]models.MovieBrief, error) {
	var allMovies []models.MovieBrief
	movieSet := make(map[string]bool) // To avoid duplicates

	// Search terms that are likely to return movies of the specified genre
	searchTerms := s.getGenreSearchTerms(genre)

	for _, term := range searchTerms {
		if len(allMovies) >= limit*2 { // Get more than needed for better filtering
			break
		}

		// Search multiple pages for each term
		for page := 1; page <= 3; page++ {
			searchResp, err := s.SearchMovies(term, page)
			if err != nil || searchResp.Response == "False" {
				continue
			}

			for _, result := range searchResp.Search {
				if movieSet[result.ImdbID] {
					continue // Skip duplicates
				}

				// Get full movie details
				movieDetails, err := s.GetMovieByTitle(result.Title)
				if err != nil || movieDetails.Response == "False" {
					continue
				}

				// Check if movie contains the desired genre
				if strings.Contains(strings.ToLower(movieDetails.Genre), strings.ToLower(genre)) {
					rating, _ := strconv.ParseFloat(movieDetails.ImdbRating, 64)
					if rating > 0 { // Only include movies with valid ratings
						movie := models.MovieBrief{
							Title:      movieDetails.Title,
							Year:       movieDetails.Year,
							ImdbRating: movieDetails.ImdbRating,
							Genre:      movieDetails.Genre,
							Director:   movieDetails.Director,
							Plot:       movieDetails.Plot,
						}
						allMovies = append(allMovies, movie)
						movieSet[result.ImdbID] = true

						if len(allMovies) >= limit*2 {
							break
						}
					}
				}
			}
		}
	}

	// Sort by IMDb rating (descending)
	sort.Slice(allMovies, func(i, j int) bool {
		ratingI, _ := strconv.ParseFloat(allMovies[i].ImdbRating, 64)
		ratingJ, _ := strconv.ParseFloat(allMovies[j].ImdbRating, 64)
		return ratingI > ratingJ
	})

	// Return top movies up to the limit
	if len(allMovies) > limit {
		allMovies = allMovies[:limit]
	}

	return allMovies, nil
}

// GetRecommendations provides movie recommendations based on a favorite movie
func (s *OMDbService) GetRecommendations(favoriteMovie string) (*models.RecommendationsResponse, error) {
	// Get details of the favorite movie
	movieDetails, err := s.GetMovieByTitle(favoriteMovie)
	if err != nil || movieDetails.Response == "False" {
		return nil, fmt.Errorf("favorite movie not found: %s", favoriteMovie)
	}

	recommendations := &models.RecommendationsResponse{
		FavoriteMovie: movieDetails.Title,
		Recommendations: models.RecommendationsByCategory{
			GenreBased:    []models.MovieBrief{},
			DirectorBased: []models.MovieBrief{},
			ActorBased:    []models.MovieBrief{},
		},
	}

	// Level 1: Genre-based recommendations
	genres := strings.Split(movieDetails.Genre, ", ")
	for _, genre := range genres {
		if len(recommendations.Recommendations.GenreBased) >= 20 {
			break
		}
		genreMovies, err := s.getMoviesExcluding(genre, "genre", movieDetails.Title, 20-len(recommendations.Recommendations.GenreBased))
		if err == nil {
			recommendations.Recommendations.GenreBased = append(recommendations.Recommendations.GenreBased, genreMovies...)
		}
	}

	// Level 2: Director-based recommendations
	directors := strings.Split(movieDetails.Director, ", ")
	for _, director := range directors {
		if len(recommendations.Recommendations.DirectorBased) >= 20 {
			break
		}
		directorMovies, err := s.getMoviesExcluding(director, "director", movieDetails.Title, 20-len(recommendations.Recommendations.DirectorBased))
		if err == nil {
			recommendations.Recommendations.DirectorBased = append(recommendations.Recommendations.DirectorBased, directorMovies...)
		}
	}

	// Level 3: Actor-based recommendations
	actors := strings.Split(movieDetails.Actors, ", ")
	for _, actor := range actors[:min(3, len(actors))] { // Limit to first 3 actors
		if len(recommendations.Recommendations.ActorBased) >= 20 {
			break
		}
		actorMovies, err := s.getMoviesExcluding(actor, "actor", movieDetails.Title, 20-len(recommendations.Recommendations.ActorBased))
		if err == nil {
			recommendations.Recommendations.ActorBased = append(recommendations.Recommendations.ActorBased, actorMovies...)
		}
	}

	// Sort each category by IMDb rating
	s.sortMoviesByRating(recommendations.Recommendations.GenreBased)
	s.sortMoviesByRating(recommendations.Recommendations.DirectorBased)
	s.sortMoviesByRating(recommendations.Recommendations.ActorBased)

	return recommendations, nil
}

// Helper function to get movies by criteria while excluding a specific movie
func (s *OMDbService) getMoviesExcluding(searchTerm, searchType, excludeTitle string, limit int) ([]models.MovieBrief, error) {
	var movies []models.MovieBrief
	movieSet := make(map[string]bool)

	// Search for movies
	for page := 1; page <= 2; page++ {
		searchResp, err := s.SearchMovies(searchTerm, page)
		if err != nil || searchResp.Response == "False" {
			continue
		}

		for _, result := range searchResp.Search {
			if movieSet[result.ImdbID] || strings.EqualFold(result.Title, excludeTitle) {
				continue
			}

			movieDetails, err := s.GetMovieByTitle(result.Title)
			if err != nil || movieDetails.Response == "False" {
				continue
			}

			// Check if movie matches the search criteria
			var matches bool
			switch searchType {
			case "genre":
				matches = strings.Contains(strings.ToLower(movieDetails.Genre), strings.ToLower(searchTerm))
			case "director":
				matches = strings.Contains(strings.ToLower(movieDetails.Director), strings.ToLower(searchTerm))
			case "actor":
				matches = strings.Contains(strings.ToLower(movieDetails.Actors), strings.ToLower(searchTerm))
			}

			if matches {
				rating, _ := strconv.ParseFloat(movieDetails.ImdbRating, 64)
				if rating > 0 {
					movie := models.MovieBrief{
						Title:      movieDetails.Title,
						Year:       movieDetails.Year,
						ImdbRating: movieDetails.ImdbRating,
						Genre:      movieDetails.Genre,
						Director:   movieDetails.Director,
						Plot:       movieDetails.Plot,
					}
					movies = append(movies, movie)
					movieSet[result.ImdbID] = true

					if len(movies) >= limit {
						return movies, nil
					}
				}
			}
		}
	}

	return movies, nil
}

// Helper function to sort movies by IMDb rating
func (s *OMDbService) sortMoviesByRating(movies []models.MovieBrief) {
	sort.Slice(movies, func(i, j int) bool {
		ratingI, _ := strconv.ParseFloat(movies[i].ImdbRating, 64)
		ratingJ, _ := strconv.ParseFloat(movies[j].ImdbRating, 64)
		return ratingI > ratingJ
	})
}

// Helper function to get search terms for different genres
func (s *OMDbService) getGenreSearchTerms(genre string) []string {
	genreTerms := map[string][]string{
		"action":    {"action", "adventure", "superhero", "martial arts", "spy"},
		"comedy":    {"comedy", "funny", "humor", "romantic comedy", "parody"},
		"drama":     {"drama", "emotional", "family", "biographical", "historical"},
		"horror":    {"horror", "scary", "thriller", "supernatural", "zombie"},
		"sci-fi":    {"science fiction", "sci-fi", "space", "future", "alien"},
		"romance":   {"romance", "love", "romantic", "wedding", "relationship"},
		"thriller":  {"thriller", "suspense", "mystery", "crime", "psychological"},
		"animation": {"animation", "animated", "cartoon", "pixar", "disney"},
		"fantasy":   {"fantasy", "magic", "wizard", "medieval", "adventure"},
		"crime":     {"crime", "gangster", "mafia", "detective", "police"},
	}

	terms, exists := genreTerms[strings.ToLower(genre)]
	if !exists {
		return []string{genre, "movie", "film"}
	}
	return terms
}

// Helper function to make HTTP requests to OMDb API
func (s *OMDbService) makeRequest(params url.Values) (*models.OMDbResponse, error) {
	resp, err := s.Client.Get(OMDbBaseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var omdbResp models.OMDbResponse
	if err := json.Unmarshal(body, &omdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if omdbResp.Response == "False" {
		return nil, fmt.Errorf("OMDb API error: %s", omdbResp.Error)
	}

	return &omdbResp, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
