package entities

import "github.com/google/uuid"

// Text entity
type Text struct {
	Id   uuid.UUID `json:"id"`
	Text string    `json:"text"`
}
