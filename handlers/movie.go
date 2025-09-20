package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"go-api/models"
	"go-api/services"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	omdbService *services.OMDbService
}

func NewMovieHandler(omdbService *services.OMDbService) *MovieHandler {
	return &MovieHandler{
		omdbService: omdbService,
	}
}

// GetMovieDetails handles GET /api/movie?title=MovieTitle
func (h *MovieHandler) GetMovieDetails(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Title query parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	movieData, err := h.omdbService.GetMovieByTitle(title)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Movie not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "Movie not found: " + title,
				Code:    http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to fetch movie details: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := models.MovieDetailsResponse{
		Title:    movieData.Title,
		Year:     movieData.Year,
		Plot:     movieData.Plot,
		Country:  movieData.Country,
		Awards:   movieData.Awards,
		Director: movieData.Director,
		Ratings:  movieData.Ratings,
	}

	c.JSON(http.StatusOK, response)
}

// GetEpisodeDetails handles GET /api/episode?series_title=SeriesTitle&season=1&episode_number=1
func (h *MovieHandler) GetEpisodeDetails(c *gin.Context) {
	seriesTitle := c.Query("series_title")
	seasonStr := c.Query("season")
	episodeStr := c.Query("episode_number")

	if seriesTitle == "" || seasonStr == "" || episodeStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "series_title, season, and episode_number query parameters are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	season, err := strconv.Atoi(seasonStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Season must be a valid integer",
			Code:    http.StatusBadRequest,
		})
		return
	}

	episode, err := strconv.Atoi(episodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Episode number must be a valid integer",
			Code:    http.StatusBadRequest,
		})
		return
	}

	episodeData, err := h.omdbService.GetEpisodeDetails(seriesTitle, season, episode)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Episode not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "Episode not found for the given parameters",
				Code:    http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to fetch episode details: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := models.EpisodeDetailsResponse{
		Title:       episodeData.Title,
		SeriesTitle: seriesTitle,
		Season:      episodeData.Season,
		Episode:     episodeData.Episode,
		Year:        episodeData.Year,
		Plot:        episodeData.Plot,
		Director:    episodeData.Director,
		Actors:      episodeData.Actors,
		Ratings:     episodeData.Ratings,
		ImdbRating:  episodeData.ImdbRating,
	}

	c.JSON(http.StatusOK, response)
}

// GetMoviesByGenre handles GET /api/movies/genre?genre=Action
func (h *MovieHandler) GetMoviesByGenre(c *gin.Context) {
	genre := c.Query("genre")
	if genre == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Genre query parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	movies, err := h.omdbService.GetMoviesByGenre(genre, 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to fetch movies by genre: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if len(movies) == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Not Found",
			Message: "No movies found for genre: " + genre,
			Code:    http.StatusNotFound,
		})
		return
	}

	response := models.GenreMoviesResponse{
		Genre:  genre,
		Movies: movies,
		Count:  len(movies),
	}

	c.JSON(http.StatusOK, response)
}

// GetRecommendations handles GET /api/recommendations?favorite_movie=MovieTitle
func (h *MovieHandler) GetRecommendations(c *gin.Context) {
	favoriteMovie := c.Query("favorite_movie")
	if favoriteMovie == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "favorite_movie query parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	recommendations, err := h.omdbService.GetRecommendations(favoriteMovie)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "Favorite movie not found: " + favoriteMovie,
				Code:    http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate recommendations: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// HealthCheck handles GET /health
func (h *MovieHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Movie API is running",
		"version": "1.0.0",
	})
}
