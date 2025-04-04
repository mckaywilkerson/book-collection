package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mckaywilkerson/book-collection/internal/db"
	"github.com/mckaywilkerson/book-collection/internal/models"
)

// HandleGetAllBooks - returns all books from the database of given pool
func HandleGetAllBooks(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := db.GetAllBooks(database)
		if err != nil {
			log.Printf("Failed to load all books from DB: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, books)
			return
		}

		c.IndentedJSON(http.StatusOK, books)
	}
}

func HandleAddBook(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bookToAdd models.Book

		err := c.BindJSON(&bookToAdd)
		if err != nil {
			log.Printf("Failed to retrieve data from body: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		id, err := db.AddBook(database, bookToAdd)
		if err != nil {
			log.Printf("Failed to add the book to database: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		bookToAdd.ID = id
		c.IndentedJSON(http.StatusCreated, bookToAdd)
	}
}

func HandleGetBook(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("Failed to convert ID to integer: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		book, err := db.GetBook(database, id)
		if err != nil {
			log.Printf("Failed to get book with this ID: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		c.IndentedJSON(http.StatusOK, book)
	}
}

func HandleUpdateBook(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("Failed to convert ID to integer: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		var newBookInfo models.Book
		err = c.BindJSON(&newBookInfo)
		if err != nil {
			log.Printf("Failed to retrieve data from body: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		err = db.UpdateBook(database, id, newBookInfo)
		if err != nil {
			log.Printf("Failed to update the book in the database: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		newBookInfo.ID = id
		c.IndentedJSON(http.StatusCreated, newBookInfo)
	}
}

func HandleDeleteBook(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("Failed to convert ID to integer: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		err = db.DeleteBook(database, id)
		if err != nil {
			log.Printf("Failed to delete book from database: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		c.IndentedJSON(http.StatusNoContent, nil)
	}
}
