package entities

import "github.com/google/uuid"

// Source entity
type Source struct {
	Id  uuid.UUID `json:"id"`
	Url string    `json:"url"`
}
