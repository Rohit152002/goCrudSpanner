package routes

import (
	"crudspanner/controller"
	"crudspanner/repositories"
	"crudspanner/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine, logger *zap.Logger, db *gorm.DB) {

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	userController := controller.NewUserController(userService)

	router.GET("/:id", userController.GetUserByID)

	router.POST("/", userController.RegistrationUser)

	router.GET("/users", userController.GetAllUsers)

	router.PUT("/:id", userController.UpdateUser)

	router.DELETE("/:id", userController.DeleteUser)
}
