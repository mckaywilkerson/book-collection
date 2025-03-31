package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mckaywilkerson/book-collection/internal/models"
)

// ConnectDB - connects to a given postgres database, via pgxpool (for simultaneous connections)
func ConnectDB(dsn string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres DB!")

	return dbPool, nil
}

// GetAllBooks - returns list of all books in the database
func GetAllBooks(dbPool *pgxpool.Pool) ([]models.Book, error) {
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	statement := "SELECT id, title, author, publication_year FROM books"

	pgxRows, err := conn.Query(context.Background(), statement)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	var books []models.Book
	for pgxRows.Next() {
		var singleBook models.Book
		if err := pgxRows.Scan(&singleBook.ID, &singleBook.Title, &singleBook.Author, &singleBook.PublicationYear); err != nil {
			return nil, err
		}
		books = append(books, singleBook)
	}
	if err := pgxRows.Err(); err != nil {
		return nil, err
	}

	return books, nil

}
