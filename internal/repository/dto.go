package repository

import "github.com/google/uuid"

// Text data transfer object entity for database requests
type Text struct {
	Id   uuid.UUID `db:"id"`
	Text string    `db:"text"`
}

// Technology data transfer object entity for database requests
type Technology struct {
	Id    uuid.UUID `db:"id"`
	Title string    `db:"title"`
	Url   string    `db:"url"`
}

// Example data transfer object entity for database requests
type Example struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}

// Source data transfer object entity for database requests
type Source struct {
	Id  uuid.UUID `db:"id"`
	Url string    `db:"url"`
}

// Lib data transfer object entity for database requests
type Lib struct {
	Id  uuid.UUID `db:"id"`
	Url string    `db:"url"`
}

// Link data transfer object entity for database requests
type Link struct {
	Id    uuid.UUID `db:"id"`
	Title string    `db:"title"`
	Url   string    `db:"url"`
}

// Software data transfer object entity for database requests
type Software struct {
	Id    uuid.UUID `db:"id"`
	Title string    `db:"title"`
	Url   string    `db:"url"`
}
