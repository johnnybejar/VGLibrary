package models

type Game struct {
	Id                    int     `json:"id"`
	AggregatedRating      float32 `json:"aggregated_rating"`
	AggregatedRatingCount int     `json:"aggregated_rating_count"`
	AlternativeNames      []int   `json:"alternative_names"`
	Artworks              []int   `json:"artworks"`
	Category              int     `json:"category"`
	Cover                 int     `json:"cover"`
	CreatedAt             int     `json:"created_at"`
	ExternalGames         []int   `json:"external_games"`
	Franchise             string  `json:"franchise"`
	FirstReleaseDate      int     `json:"first_release_date"`
	GameEngines           []int   `json:"game_engines"`
	GameModes             []int   `json:"game_modes"`
	Genres                []int   `json:"genres"`
	InvolvedCompanies     []int   `json:"involved_companies"`
	Name                  string  `json:"name"`
	Platforms             []int   `json:"platforms"`
	ReleaseDates          []int   `json:"release_dates"`
	Slug                  string  `json:"slug"`
	Summary               string  `json:"summary"`
	UpdatedAt             int     `json:"updated_at"`
	Url                   string  `json:"url"`
	Websites              []int   `json:"websites"`
	Checksum              string  `json:"checksum"`
}