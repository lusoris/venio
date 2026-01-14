package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func seedTestUsers(ctx context.Context, pool *pgxpool.Pool) error {
	testUsers := []struct {
		email    string
		password string
		roleName string
	}{
		{
			email:    "admin@test.local",
			password: "AdminPassword123!",
			roleName: "admin",
		},
		{
			email:    "moderator@test.local",
			password: "ModeratorPassword123!",
			roleName: "moderator",
		},
		{
			email:    "user@test.local",
			password: "UserPassword123!",
			roleName: "user",
		},
		{
			email:    "guest@test.local",
			password: "GuestPassword123!",
			roleName: "guest",
		},
	}

	for _, testUser := range testUsers {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		// Generate verification token
		token := generateToken()

		// Insert user
		query := `
			INSERT INTO users (email, password, verification_token, is_verified, created_at)
			VALUES ($1, $2, $3, true, NOW())
			ON CONFLICT (email) DO NOTHING
			RETURNING id
		`

		var userID int64
		err = pool.QueryRow(ctx, query, testUser.email, string(hashedPassword), token).Scan(&userID)
		if err != nil {
			log.Printf("Warning: Could not insert user %s: %v", testUser.email, err)
			// Try to get existing user ID
			getQuery := `SELECT id FROM users WHERE email = $1`
			err = pool.QueryRow(ctx, getQuery, testUser.email).Scan(&userID)
			if err != nil {
				return fmt.Errorf("get user id for %s: %w", testUser.email, err)
			}
		}

		// Assign role
		roleQuery := `
			INSERT INTO user_roles (user_id, role_id, created_at)
			SELECT $1, id, NOW()
			FROM roles
			WHERE name = $2
			ON CONFLICT (user_id, role_id) DO NOTHING
		`

		_, err = pool.Exec(ctx, roleQuery, userID, testUser.roleName)
		if err != nil {
			return fmt.Errorf("assign role %s to user %s: %w", testUser.roleName, testUser.email, err)
		}

		log.Printf("✓ Test user created: %s (role: %s)", testUser.email, testUser.roleName)
	}

	return nil
}

func generateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://venio:venio@localhost:5432/venio?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Create connection pool: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Connect to database: %v", err)
	}

	log.Println("Connected to database. Seeding test users...")

	if err := seedTestUsers(ctx, pool); err != nil {
		log.Fatalf("Seed test users: %v", err)
	}

	log.Println("✓ Test user seeding complete!")
	log.Println("")
	log.Println("Test User Credentials:")
	log.Println("  Admin:      admin@test.local / AdminPassword123!")
	log.Println("  Moderator:  moderator@test.local / ModeratorPassword123!")
	log.Println("  User:       user@test.local / UserPassword123!")
	log.Println("  Guest:      guest@test.local / GuestPassword123!")
}
