package main

import (
	"crudspanner/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"crudspanner/config"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic("Failed to initialize logger" + err.Error())
	}
	defer logger.Sync()
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load env file")
	}

	httpServer := gin.Default()
	dbConnector:= config.ConnectDB()
	routes.UserRoutes(httpServer, logger,dbConnector)
	httpServer.Run(":8080")

}
