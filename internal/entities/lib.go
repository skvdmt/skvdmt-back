package entities

import "github.com/google/uuid"

// Lib entity
type Lib struct {
	Id  uuid.UUID `json:"id"`
	Url string    `json:"url"`
}
