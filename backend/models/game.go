package models

type Game struct {
	Id                    int      `json:"id"`
	AggregatedRating      float32  `json:"aggregated_rating"`
	AggregatedRatingCount int      `json:"aggregated_rating_count"`
	Collections           []int    `json:"collections"` // Refers to franchise of the game
	Cover                 int      `json:"cover"`
	GameModes             []int    `json:"game_modes"`
	GameModesString       []string // Declared after added to list
	Genres                []int    `json:"genres"`
	GenresString          []string // Declared after added to list
	InvolvedCompanies     []int    `json:"involved_companies"`
	Name                  string   `json:"name"`
	Platforms             []int    `json:"platforms"`
	ReleaseDate           int      `json:"first_release_date"`
	Slug                  string   `json:"slug"`
	Summary               string   `json:"summary"`
	Url                   string   `json:"url"`
}

type Collections struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Cover struct {
	Id      int    `json:"id"`
	Game    int    `json:"game"`
	ImageId string `json:"image_id"`
}
