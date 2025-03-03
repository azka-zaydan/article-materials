package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database configuration
const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "test"
)

// DB instance
var db *sqlx.DB

func initDB() error {
	// PostgreSQL connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	log.Println("Connected to PostgreSQL successfully!")
	return nil
}
