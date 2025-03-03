package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	multipleUserCreate()

	ctx := context.Background()
	// get all users and tokens
	users, err := GetAllUserAndTokens(ctx)
	if err != nil {
		log.Fatalf("Failed to get all users and tokens: %v", err)
	}

	for _, u := range users {
		fmt.Printf("User: %s, Email: %s\n", u.Name, u.Email)
	}

}

func multipleUserCreate() {
	var wg sync.WaitGroup
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			userID, err := uuid.NewV4()
			if err != nil {
				log.Fatalf("Failed to generate UUID: %v", err)
			}

			user := User{
				ID:    userID.String(),
				Name:  generateRandomName(),
				Email: generateRandomEmail(),
			}

			err = CreateUserWithToken(ctx, user)
			if err != nil {
				log.Fatalf("Failed to create user with token: %v", err)
			}
		}()
	}

	wg.Wait()
	fmt.Println("All users and tokens created successfully.")

}

func CreateUserWithToken(ctx context.Context, user User) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = CreateUser(ctx, tx, &user)
	if err != nil {
		return err
	}

	token := generateToken(user.ID)
	err = CreateUserToken(ctx, tx, user.ID, token)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CreateUser(ctx context.Context, tx *sqlx.Tx, user *User) error {
	query := "INSERT INTO users (id, name, email, created_at) VALUES (:id, :name, :email, NOW())"
	_, err := tx.ExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return err
}

func CreateUserToken(ctx context.Context, tx *sqlx.Tx, userID, token string) error {
	query := "INSERT INTO user_tokens (user_id, token, created_at) VALUES ($1, $2, NOW())"
	_, err := tx.ExecContext(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("failed to create user token: %w", err)
	}
	return err
}

func generateToken(userID string) string {
	return fmt.Sprintf("token-%s", userID)
}

func GetAllUserAndTokens(ctx context.Context) ([]User, error) {
	query := `
		SELECT u.id, u.name, u.email
		FROM users u
		JOIN user_tokens ut ON u.id = ut.user_id
	`
	var users []User
	err := db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

func generateRandomName() string {
	names := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hannah"}
	return names[rand.Intn(len(names))]
}

func generateRandomEmail() string {
	domains := []string{"example.com", "test.com", "mail.com", "random.org"}
	name := generateRandomName()
	return fmt.Sprintf("%s@%s", name, domains[rand.Intn(len(domains))])
}
