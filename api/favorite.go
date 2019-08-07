package api

import "time"

type FavoritePage struct {
	Data       []Favorite `json:"data"`
	NextCursor string     `json:"nextCursor,omitempty"`
}

type Favorite struct {
	ID           string    `json:"id"`
	Owner        Identity  `json:"owner"`
	Created      time.Time `json:"created"`
	Description  string    `json:"description,omitempty"`
	Name         string    `json:"name,omitempty"`
	FavoriteType string    `json:"type"`
}
