package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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
