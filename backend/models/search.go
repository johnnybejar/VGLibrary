package models

// Will use the same endpoint as the Game type, but with reduced fields to minimize data
type Search struct {
	Id int `json:"id"`
	Artworks int `json:"artworks"` 
	Cover []int `json:"cover"`
	Franchise string `json:"franchise"`
	Genres []int `json:"genres"`
	Name string `json:"name"`
}