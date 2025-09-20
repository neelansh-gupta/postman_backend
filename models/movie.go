package models

// OMDbResponse represents the full response from OMDb API
type OMDbResponse struct {
	Title      string   `json:"Title"`
	Year       string   `json:"Year"`
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Runtime    string   `json:"Runtime"`
	Genre      string   `json:"Genre"`
	Director   string   `json:"Director"`
	Writer     string   `json:"Writer"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	ImdbRating string   `json:"imdbRating"`
	ImdbVotes  string   `json:"imdbVotes"`
	ImdbID     string   `json:"imdbID"`
	Type       string   `json:"Type"`
	DVD        string   `json:"DVD"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`
	Response   string   `json:"Response"`
	Error      string   `json:"Error,omitempty"`
	Season     string   `json:"Season,omitempty"`
	Episode    string   `json:"Episode,omitempty"`
}

// Rating represents a rating from different sources
type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

// MovieDetailsResponse represents the cleaned response for movie details
type MovieDetailsResponse struct {
	Title    string   `json:"title"`
	Year     string   `json:"year"`
	Plot     string   `json:"plot"`
	Country  string   `json:"country"`
	Awards   string   `json:"awards"`
	Director string   `json:"director"`
	Ratings  []Rating `json:"ratings"`
}

// EpisodeDetailsResponse represents the cleaned response for episode details
type EpisodeDetailsResponse struct {
	Title       string   `json:"title"`
	SeriesTitle string   `json:"series_title"`
	Season      string   `json:"season"`
	Episode     string   `json:"episode"`
	Year        string   `json:"year"`
	Plot        string   `json:"plot"`
	Director    string   `json:"director"`
	Actors      string   `json:"actors"`
	Ratings     []Rating `json:"ratings"`
	ImdbRating  string   `json:"imdb_rating"`
}

// GenreMoviesResponse represents the response for genre-based movies
type GenreMoviesResponse struct {
	Genre  string       `json:"genre"`
	Movies []MovieBrief `json:"movies"`
	Count  int          `json:"count"`
}

// MovieBrief represents a brief movie information
type MovieBrief struct {
	Title      string `json:"title"`
	Year       string `json:"year"`
	ImdbRating string `json:"imdb_rating"`
	Genre      string `json:"genre"`
	Director   string `json:"director"`
	Plot       string `json:"plot"`
}

// RecommendationsResponse represents movie recommendations
type RecommendationsResponse struct {
	FavoriteMovie   string                    `json:"favorite_movie"`
	Recommendations RecommendationsByCategory `json:"recommendations"`
}

// RecommendationsByCategory categorizes recommendations by priority
type RecommendationsByCategory struct {
	GenreBased    []MovieBrief `json:"genre_based"`
	DirectorBased []MovieBrief `json:"director_based"`
	ActorBased    []MovieBrief `json:"actor_based"`
}

// SearchResponse represents OMDb search results
type SearchResponse struct {
	Search       []SearchResult `json:"Search"`
	TotalResults string         `json:"totalResults"`
	Response     string         `json:"Response"`
	Error        string         `json:"Error,omitempty"`
}

// SearchResult represents individual search result
type SearchResult struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
