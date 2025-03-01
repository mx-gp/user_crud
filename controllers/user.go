package controllers

import (
	"net/http"
	"strconv"
	"user_crud/models"
	"user_crud/repository"
	"user_crud/utils"

	"github.com/gin-gonic/gin"
)

// Create User
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := repository.CreateUser(user)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SendSuccessResponse(c, "User created successfully", nil)
}

// Get All Users
func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUsers()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	utils.SendSuccessResponse(c, "Users fetched successfully", users)
}

// Get User by ID
func GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := repository.GetUserByID(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	utils.SendSuccessResponse(c, "User fetched successfully", user)
}

// Update User
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := repository.UpdateUser(user)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.SendSuccessResponse(c, "User updated successfully", nil)
}

// Delete User
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := repository.DeleteUser(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.SendSuccessResponse(c, "User deleted successfully", nil)
}
