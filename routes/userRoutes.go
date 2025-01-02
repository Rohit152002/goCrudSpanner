package routes

import (
	"crudspanner/controller"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine, logger *zap.Logger, db *gorm.DB) {
	userController := controller.UserController{Logger: logger, Db: db}

	router.GET("/:id", userController.GetUserById)
	router.POST("/", userController.CreateUser)
	router.GET("/users", userController.GetAllUser)
	router.PUT("/:id", userController.UpdateUser)
	router.DELETE("/:id", userController.DeleteUser)
}
