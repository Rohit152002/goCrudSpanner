package controller

import (
	"crudspanner/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserController struct {
	Logger *zap.Logger
	Db     *gorm.DB
}

func (uc *UserController) GetUserById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}
	var user model.User

	if err := uc.Db.Where("id=? ", id).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAllUser(ctx *gin.Context) {
	var users []model.User
	uc.Db.Find(&users)
	if users == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "There is no user"})
	}
	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}
	var user model.User
	if err := uc.Db.Where("id=? ", id).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	uc.Db.Delete(&user)
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}
	var user model.User
	if err := uc.Db.Where("id=? ", id).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var updatedUser model.User
	ctx.ShouldBindJSON(&updatedUser)
	user.Name = updatedUser.Name
	user.Email = updatedUser.Email
	user.Address = updatedUser.Address
	user.Password = updatedUser.Password
	if err := uc.Db.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	ctx.ShouldBindJSON(&user)

	var existingUser model.User

	userNotFoundError := uc.Db.Where("email = ?", user.Email).First(&existingUser).Error

	if userNotFoundError == gorm.ErrRecordNotFound {

		newUser := &model.User{Name: user.Name, Email: user.Email, Password: user.Password, Address: user.Address}

		primaryKey := uc.Db.Create(newUser)

		if primaryKey.Error != nil {
			uc.Logger.Error("Failed to Create user", zap.String("userName ", user.Name), zap.Error(primaryKey.Error))
			ctx.JSON(http.StatusConflict, gin.H{"message": "The Phone is already registered"})
			return
		}
		uc.Logger.Info(fmt.Sprintf("User %s created successfully", user.Name))

		ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	} else {
		uc.Logger.Warn("User Email Already Exist", zap.String("usermail", user.Email))
		ctx.JSON(http.StatusConflict, gin.H{"message": "User Email Already Exist"})
	}
}
