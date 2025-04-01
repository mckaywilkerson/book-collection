package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver for use with database/sql
	"github.com/mckaywilkerson/book-collection/internal/models"
)

// ConnectDB - connects to a given postgres database, via pgxpool (for simultaneous connections)
func ConnectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres DB!")

	return db, nil
}

// GetAllBooks - returns list of all books in the database
func GetAllBooks(db *sql.DB) ([]models.Book, error) {
	statement := "SELECT * FROM books"

	dbRows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer dbRows.Close()

	var books []models.Book
	for dbRows.Next() {
		var singleBook models.Book
		if err := dbRows.Scan(&singleBook.ID, &singleBook.Title, &singleBook.Author, &singleBook.PublicationYear); err != nil {
			return nil, err
		}
		books = append(books, singleBook)
	}
	if err := dbRows.Err(); err != nil {
		return nil, err
	}

	return books, nil

}

// GetBook - returns single book information given ID
func GetBook(db *sql.DB, bookID int) (models.Book, error) {
	statement := "SELECT * FROM books WHERE id=$1"

	row := db.QueryRow(statement, bookID)
	var book models.Book

	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.PublicationYear); err != nil {
		if err == sql.ErrNoRows {
			return book, fmt.Errorf("GetBook, no such book", bookID)
		}
		return book, err
	}

	return book, nil
}

// UpdateBook - updates the book data
func UpdateBook(db *sql.DB, bookID int, newBookInfo models.Book) error {
	statement := "UPDATE books SET title=$1, author=$2, publication_year=$3 WHERE id=$4"

	if _, err := db.Exec(statement, newBookInfo.Title, newBookInfo.Author, newBookInfo.PublicationYear, bookID); err != nil {
		return err
	}

	return nil
}

// DeleteBook - delete the given book
func DeleteBook(db *sql.DB, bookID int) error {
	statement := "DELETE FROM books WHERE id=$1"

	if _, err := db.Exec(statement, bookID); err != nil {
		return err
	}

	return nil
}

// AddBook - add a new book to the database
func AddBook(db *sql.DB, bookInfo models.Book) (int64, error) {
	statement := "INSERT INTO books (title, author, publication_year) VALUES ($1, $2, $3)"

	result, err := db.Exec(statement, bookInfo.Title, bookInfo.Author, bookInfo.PublicationYear)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
