package controller

import (
	"crudspanner/model"
	"crudspanner/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// RegistrationUser godoc
// @Summary Register a new user
// @Description Create a new user in the system
// @Tags users
// @Param user body model.User true "User Data"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string "Invalid input"
// @Router / [post]
func (ctrl *UserController) RegistrationUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid input"})
		return
	}

	createdUser, err := ctrl.userService.Registration(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user by their unique ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]string "User not found"
// @Router /{id} [get]
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is missing"})
		return
	}
	user, err := ctrl.userService.GetUserById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)

}

// GetAllUsers godoc
// @Summary List all users
// @Description Retrieve all users in the system
// @Tags users
// @Success 200 {array} model.User
// @Router /users [get]
func (ctrl *UserController) GetAllUsers(c *gin.Context) {

	users, err := ctrl.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return

	}
	c.JSON(http.StatusOK, users)

}

// UpdateUser godoc
// @Summary Update user information
// @Description Update the details of an existing user
// @Tags users
// @Param id path int true "User ID"
// @Param user body model.User true "Updated User Data"
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]string "User not found"
// @Router /{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {

	var user model.User
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is missing"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	updatedUser, err := ctrl.userService.UpdateUser(uint(id), &user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)

}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user from the system by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 404 {object} map[string]string "User not found"
// @Router /{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is missing"})
		return
	}

	if err := ctrl.userService.DeleteUser(uint(id)); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})

		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})

}
