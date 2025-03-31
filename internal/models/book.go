package models

// Book is the data model for a book in the DB
type Book struct {
	ID              int    `json:"id,omitempty"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
}
