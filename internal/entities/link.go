package entities

import "github.com/google/uuid"

// Link entity
type Link struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Url   string    `json:"url"`
}
