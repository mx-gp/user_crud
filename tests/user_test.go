package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"user_crud/config"
	"user_crud/models"
	"user_crud/repository"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func setupTestDB() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir := strings.Replace(currentDir, "/models", "", 1)

	// Load environment variables
	err = godotenv.Load(dir + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	fmt.Printf("%+v", dsn)
	config.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// Create a test table (optional)
	_, err = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			email VARCHAR(100) UNIQUE,
			age INT
		);
	`)
	// fmt.Println(r.RowsAffected())
	if err != nil {
		panic(err)
	}
}

func teardownTestDB() {
	// Drop test data after each test
	config.DB.Exec("DELETE FROM users")
	config.DB.Close()
}

func TestCreateUser(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	user := models.User{Name: "Alice", Email: "alice@example.com", Age: 28}
	err := repository.CreateUser(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("User was not inserted into the database")
	}

	// Delete user
	err = repository.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}
}

func TestGetAllUsers(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert test users
	config.DB.Exec("INSERT INTO users (name, email, age) VALUES ('Bob', 'bob@example.com', 25), ('Charlie', 'charlie@example.com', 29)")

	users, err := repository.GetAllUsers()
	if err != nil {
		t.Errorf("Failed to fetch users: %v", err)
	}

	if len(users) < 2 {
		t.Errorf("Expected at least 2 users, got %d", len(users))
	}

	for _, user := range users {
		// Delete user
		err = repository.DeleteUser(user.ID)
		if err != nil {
			t.Errorf("Failed to delete user: %v", err)
		}
	}
}

func TestGetUserByID(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert a test user
	var id int
	config.DB.QueryRow("INSERT INTO users (name, email, age) VALUES ('David', 'david@example.com', 40) RETURNING id").Scan(&id)

	user, err := repository.GetUserByID(id)
	if err != nil {
		t.Errorf("Failed to get user by ID: %v", err)
	}

	if user.Email != "david@example.com" {
		t.Errorf("Expected email 'david@example.com', got %s", user.Email)
	}
	// Delete user
	err = repository.DeleteUser(id)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert test user
	var id int
	config.DB.QueryRow("INSERT INTO users (name, email, age) VALUES ('Eve', 'eve@example.com', 22) RETURNING id").Scan(&id)

	// Update the user
	updatedUser := models.User{ID: id, Name: "Eve Adams", Email: "eve@example.com", Age: 30}
	err := repository.UpdateUser(updatedUser)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	// Verify update
	var name string
	var age int
	config.DB.QueryRow("SELECT name, age FROM users WHERE id = $1", id).Scan(&name, &age)

	if name != "Eve Adams" || age != 30 {
		t.Errorf("User was not updated correctly: got name=%s, age=%d", name, age)
	}

	// Delete user
	err = repository.DeleteUser(id)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert test user
	var id int
	config.DB.QueryRow("INSERT INTO users (name, email, age) VALUES ('Frank', 'frank@example.com', 45) RETURNING id").Scan(&id)

	// Delete user
	err := repository.DeleteUser(id)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}

	// Verify deletion
	var count int
	config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&count)

	if count != 0 {
		t.Errorf("User was not deleted")
	}
}
