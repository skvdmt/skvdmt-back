package entities

import "github.com/google/uuid"

// Example entity
type Example struct {
	Id           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Links        *[]Link       `json:"links,omitempty"`
	Technologies *[]Technology `json:"technologies"`
	Sources      *[]Source     `json:"sources"`
}
