package controller

import (
	"fmt"
	"net/http"
	"strconv"

	service "github.com/dwarowski/medods-test-task/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/users/:id", func(ctx *gin.Context) { GetUserHandler(ctx, db) })
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

	user, err := service.GetByID(db, idn)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
