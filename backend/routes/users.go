package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/users"
)

func validate(c *gin.Context) {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx := c.Request.Context()
	_, err := users.GetUser(ctx, username)

	if err == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func signup(c *gin.Context) {
	var request users.SignupRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	_, err := users.GetUser(ctx, request.Username)

	if err == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	passwordHash, err := users.GenerateHash(request.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	newUser := users.UserInfo{
		Username:     request.Username,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	if err := users.CreateUser(ctx, newUser); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

func login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	user, err := users.GetUser(ctx, request.Username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := users.VerifyPassword(request.Password, user.PasswordHash); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	loginDuration := time.Now().Add(2 * time.Hour)
	token, err := auth.GenerateToken(user, loginDuration)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"role":     user.Role,
		"token":    token,
		"expires":  loginDuration.Unix(),
	})
}

func AddUserRoutes(router gin.IRouter) {
	router.GET("/users/validate/:username", validate)
	router.POST("/signup", signup)
	router.POST("/login", login)
}
