package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dwarowski/medods-test-task/src/dto"
	"github.com/dwarowski/medods-test-task/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/users/:id", func(ctx *gin.Context) { GetUserHandler(ctx, db) })
	router.POST("/register", func(ctx *gin.Context) { CreateUserHandler(ctx, db) })
}

// @Summary Get a user by ID
// @ID get-user-by-id
// @Produce json
// @Param id path int true "User ID"
// @Router /users/{id} [get]
func GetUserHandler(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")

	idn, err := strconv.Atoi(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		fmt.Println(idn)
	}

	user, err := services.GetByID(db, idn)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary create a user
// @ID create-user
// @Produce json
// @Param dto body dto.CreateUserDto true "register user"
// @Router /register [post]
func CreateUserHandler(ctx *gin.Context, db *gorm.DB) {
	var dto dto.CreateUserDto
	ctx.BindJSON(&dto)

	result, err := services.CreateUser(db, dto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
