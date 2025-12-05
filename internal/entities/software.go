package entities

import "github.com/google/uuid"

// Software entity
type Software struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Url   string    `json:"url"`
}
