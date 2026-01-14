package entities

import "github.com/google/uuid"

// Text Текст.
type Text struct {
	Id   uuid.UUID `json:"id"`
	Text string    `json:"text"`
}
