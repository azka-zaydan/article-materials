package infras

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
var DB *sqlx.DB

func InitDB() error {
	// PostgreSQL connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection settings
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Minute * 5)

	log.Println("Connected to PostgreSQL successfully!")
	return nil
}
