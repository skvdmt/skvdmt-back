package entities

import "github.com/google/uuid"

// Technology Технология.
type Technology struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Url   string    `json:"url"`
}
