package models

import (
	"database/sql"
	"testing"

	"user_crud/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func setupTestDB() {
	// Initialize the test database connection
	var err error
	config.DB, err = sql.Open("postgres", "host=127.0.0.1 port=5432 user=mx dbname=crud_db sslmode=disable")
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

	user := User{Name: "Alice", Email: "alice@example.com", Age: 28}
	err := CreateUser(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("User was not inserted into the database")
	}
}

func TestGetAllUsers(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert test users
	config.DB.Exec("INSERT INTO users (name, email, age) VALUES ('Bob', 'bob@example.com', 25), ('Charlie', 'charlie@example.com', 29)")

	users, err := GetAllUsers()
	if err != nil {
		t.Errorf("Failed to fetch users: %v", err)
	}

	if len(users) < 2 {
		t.Errorf("Expected at least 2 users, got %d", len(users))
	}
}

func TestGetUserByID(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert a test user
	var id int
	config.DB.QueryRow("INSERT INTO users (name, email, age) VALUES ('David', 'david@example.com', 40) RETURNING id").Scan(&id)

	user, err := GetUserByID(id)
	if err != nil {
		t.Errorf("Failed to get user by ID: %v", err)
	}

	if user.Email != "david@example.com" {
		t.Errorf("Expected email 'david@example.com', got %s", user.Email)
	}
}

func TestUpdateUser(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert test user
	var id int
	config.DB.QueryRow("INSERT INTO users (name, email, age) VALUES ('Eve', 'eve@example.com', 22) RETURNING id").Scan(&id)

	// Update the user
	updatedUser := User{ID: id, Name: "Eve Adams", Email: "eve@example.com", Age: 30}
	err := UpdateUser(updatedUser)
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
}

func TestDeleteUser(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// Insert test user
	var id int
	config.DB.QueryRow("INSERT INTO users (name, email, age) VALUES ('Frank', 'frank@example.com', 45) RETURNING id").Scan(&id)

	// Delete user
	err := DeleteUser(id)
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
