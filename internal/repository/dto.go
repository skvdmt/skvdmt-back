package repository

import "github.com/google/uuid"

// Text Объект передачи данных для запросов к базе данных.
type Text struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Text string    `db:"text"`
}

// Technology Объект передачи данных для запросов к базе данных.
type Technology struct {
	Id    uuid.UUID `db:"id"`
	Title string    `db:"title"`
	Url   string    `db:"url"`
}

// Example Объект передачи данных для запросов к базе данных.
type Example struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}

// Source Объект передачи данных для запросов к базе данных.
type Source struct {
	Id  uuid.UUID `db:"id"`
	Url string    `db:"url"`
}

// Lib Объект передачи данных для запросов к базе данных.
type Lib struct {
	Id  uuid.UUID `db:"id"`
	Url string    `db:"url"`
}

// Link Объект передачи данных для запросов к базе данных.
type Link struct {
	Id    uuid.UUID `db:"id"`
	Title string    `db:"title"`
	Url   string    `db:"url"`
}

// Software Объект передачи данных для запросов к базе данных.
type Software struct {
	Id    uuid.UUID `db:"id"`
	Title string    `db:"title"`
	Url   string    `db:"url"`
}
