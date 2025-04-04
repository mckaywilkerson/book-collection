package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mckaywilkerson/book-collection/internal/models"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	_, thisFile, _, _ := runtime.Caller(0)
	base := filepath.Dir(thisFile)
	relPath := "../../deploy/docker/init.sql"
	absPath, err := filepath.Abs(filepath.Join(base, relPath))
	if err != nil {
		log.Panicf("Absolute path is wrong: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=books",
			"listen_addresses = '*'",
		},
		// Mount init.sql file into /docker-entrypoint-initdb.d
		Mounts: []string{
			fmt.Sprintf("%s:/docker-entrypoint-initdb.d/init.sql", absPath),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://user_name:secret@%s/books?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url:", databaseURL)

	resource.Expire(120)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	err = pool.Retry(func() error {
		db, err = sql.Open("pgx", databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func TestGetAllBooks(t *testing.T) {
	myBooks, err := GetAllBooks(db)

	if err != nil {
		t.Error("ran into error with GetAllBooks:", err)
	}

	if len(myBooks) != 0 {
		t.Error("Expected empty slice, got", myBooks)
	}
}

func TestAddBook(t *testing.T) {
	book1 := models.Book{
		Title:           "test book 1",
		Author:          "author one",
		PublicationYear: 2000,
	}

	id, err := AddBook(db, book1)
	if err != nil {
		t.Error("Error adding book 1 during testing:", err)
	}

	book1.ID = id
	bookVerification, err := GetBook(db, id)
	if err != nil {
		t.Error("Unable to find book 1 during testing:", err)
	}

	if book1 != bookVerification {
		t.Error("Add failed: book1-", book1, "bookverification-", bookVerification)
	}
}

func TestUpdateBook(t *testing.T) {
	newBookInfo := models.Book{
		Title:           "newTitle",
		Author:          "newAuthor",
		PublicationYear: 2012,
	}

	books, err := GetAllBooks(db)
	if err != nil {
		t.Error("Error getting books:", err)
	}

	bookID := books[0].ID
	newBookInfo.ID = bookID

	err = UpdateBook(db, bookID, newBookInfo)
	if err != nil {
		t.Error("Error updating book:", err)
	}

	updatedBook, err := GetBook(db, bookID)
	if err != nil {
		t.Error("Error getting book:", err)
	}

	if updatedBook != newBookInfo {
		t.Error("Update not as expected. Got", updatedBook, "Expected", newBookInfo)
	}
}

func TestDeleteBook(t *testing.T) {
	books, err := GetAllBooks(db)
	if err != nil {
		t.Error("Error getting books:", err)
	}

	if len(books) != 1 {
		t.Error("Expected one book in db, but got", len(books))
	}

	bookID := books[0].ID

	err = DeleteBook(db, bookID)
	if err != nil {
		t.Error("Error deleting book:", err)
	}

	books, err = GetAllBooks(db)
	if err != nil {
		t.Error("Error getting books after delete:", err)
	}

	if len(books) != 0 {
		t.Error("Should have 0 books after delete, but have", len(books))
	}
}

func TestAddBook_second(t *testing.T) {
	booksToAdd := []models.Book{
		{
			Title:           "harry potter 1",
			Author:          "JK Rowling",
			PublicationYear: 2000,
		},
		{
			Title:           "harry potter 2",
			Author:          "JK Rowling",
			PublicationYear: 2002,
		},
		{
			Title:           "harry potter 3",
			Author:          "JK Rowling",
			PublicationYear: 2004,
		},
		{
			Title:           "Lord of the Rings",
			Author:          "Tolkien",
			PublicationYear: 2007,
		},
	}

	for i, book := range booksToAdd {
		id, err := AddBook(db, book)
		if err != nil {
			t.Error("Error adding book:", err)
		}
		booksToAdd[i].ID = id
	}

	getBooks, err := GetAllBooks(db)
	if err != nil {
		t.Error("Error getting books after add:", err)
	}

	if !slices.Equal(getBooks, booksToAdd) {
		t.Error("Expected:", booksToAdd, "Got:", getBooks)
	}
}

func TestDeleteBook_second(t *testing.T) {
	books, err := GetAllBooks(db)
	if err != nil {
		t.Error("Error getting books:", err)
	}
	sizeBeforeDelete := len(books)
	if sizeBeforeDelete != 4 {
		t.Error("Expected db to have 4 books, but instead it had", sizeBeforeDelete)
	}

	var idsToDelete []int
	idsToDelete = append(idsToDelete, books[1].ID, books[2].ID)

	for _, id := range idsToDelete {
		err = DeleteBook(db, id)
		if err != nil {
			t.Error("Error deleting book:", id, err)
		}
	}

	booksAfterDelete, err := GetAllBooks(db)
	if err != nil {
		t.Error("Error getting books after delete calls", err)
	}

	if len(booksAfterDelete) != 2 {
		t.Error("Length does not match what is expected, Expected 2, Got", len(booksAfterDelete))
	}

	books = slices.Delete(books, 1, 3)
	if !slices.Equal(books, booksAfterDelete) {
		t.Error("Expected", books, "Got", booksAfterDelete)
	}
}
