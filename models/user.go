package models

import "user_crud/config"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Create User
func CreateUser(user User) error {
	_, err := config.DB.Exec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3)", user.Name, user.Email, user.Age)
	return err
}

// Get All Users
func GetAllUsers() ([]User, error) {
	rows, err := config.DB.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Get User by ID
func GetUserByID(id int) (User, error) {
	var user User
	err := config.DB.QueryRow("SELECT id, name, email, age FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
	return user, err
}

// Update User
func UpdateUser(user User) error {
	_, err := config.DB.Exec("UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4", user.Name, user.Email, user.Age, user.ID)
	return err
}

// Delete User
func DeleteUser(id int) error {
	_, err := config.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
