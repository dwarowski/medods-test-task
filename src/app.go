package app

import (
	_ "github.com/dwarowski/medods-test-task/docs"
	"github.com/dwarowski/medods-test-task/src/config"
	controller "github.com/dwarowski/medods-test-task/src/controllers"
	models "github.com/dwarowski/medods-test-task/src/model"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Medods Test Task API
// @version 1.0

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func SetupApp() (*gin.Engine, error) {
	db, err := config.ConnectToDatabase()
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{})

	r := gin.Default()

	controller.RegisterRoutes(r, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r, nil
}
