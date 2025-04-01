package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mckaywilkerson/book-collection/internal/db"
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
