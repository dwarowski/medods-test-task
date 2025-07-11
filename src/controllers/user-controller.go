package controllers

import (
	"errors"
	"net/http"

	"github.com/dwarowski/medods-test-task/src/dto"
	"github.com/dwarowski/medods-test-task/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/users/:id", func(ctx *gin.Context) { GetUserHandler(ctx, db) })
	router.POST("/register", func(ctx *gin.Context) { CreateUserHandler(ctx, db) })
	router.POST("/login", func(ctx *gin.Context) { LoginHandler(ctx, db) })
	router.POST("/refresh", func(ctx *gin.Context) { RefreshHandler(ctx, db) })
	router.GET("/getUUID", func(ctx *gin.Context) { getUUID(ctx, db) })
}

// @Summary Get a user by ID
// @ID get-user-by-id
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Router /users/{id} [get]
func GetUserHandler(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")

	idn, err := uuid.Parse(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.New("User-Agent Not Found")})
	}

	ipAdress := ctx.ClientIP()

	user, err := services.GetByID(db, idn, userAgent, ipAdress)
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

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.New("User-Agent Not Found")})
	}

	ipAdress := ctx.ClientIP()

	result, err := services.CreateUser(db, dto, userAgent, ipAdress)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Summary Log in
// @ID login
// @Produce json
// @Param dto body dto.LoginDto true "login user"
// @Router /login [post]
func LoginHandler(ctx *gin.Context, db *gorm.DB) {
	var dto dto.LoginDto
	ctx.BindJSON(&dto)

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.New("User-Agent Not Found")})
	}

	ipAdress := ctx.ClientIP()

	result, err := services.Login(db, dto, userAgent, ipAdress)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Summary Refresh
// @ID token-refresh
// @Produce json
// @Param dto body dto.TokensDto true "refresh tokens"
// @Router /refresh [post]
func RefreshHandler(ctx *gin.Context, db *gorm.DB) {
	var dto dto.TokensDto
	ctx.BindJSON(&dto)

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.New("User-Agent Not Found")})
	}

	ipAdress := ctx.ClientIP()

	result, err := services.RefreshToken(db, dto, userAgent, ipAdress)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Summary Get user Id
// @ID get-uuid
// @Produce json
// @Router /getUUID [get]
// @Security ApiKeyAuth
func getUUID(ctx *gin.Context, db *gorm.DB) {
	authToken := ctx.GetHeader("Authorization")
	if authToken == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "authorization Not Found"})
		return
	}

	result, err := services.GetUUID(db, authToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)

}
