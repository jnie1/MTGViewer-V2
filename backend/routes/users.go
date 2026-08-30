package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/users"
)

func signup(c *gin.Context) {
	var request users.SignupRequest

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, err := users.GetUser(request.Username); !errors.Is(err, sql.ErrNoRows) {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	passwordHash, err := users.GenerateHash(request.Password)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newUser := users.UserInfo{
		Username:     request.Username,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	if err := users.CreateUser(newUser); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func login(c *gin.Context) {
	var request users.LoginRequest

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := users.GetUser(request.Username)
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

	c.SetCookie("token", token, int(loginDuration.Unix()), "", "", false, true)
	c.Status(http.StatusNoContent)
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.Status(http.StatusNoContent)
}

func AddUserRoutes(router gin.IRouter) {
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.POST("/logout", logout)
}
