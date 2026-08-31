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

func validate(c *gin.Context) {
	var param struct {
		Username string `uri:"username" binding:"required"`
	}

	if err := c.BindUri(&param); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx := c.Request.Context()
	_, err := users.GetUser(ctx, param.Username)

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
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	_, err := users.GetUser(ctx, req.Username)

	if err == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	passwordHash, err := users.GenerateHash(req.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	newUser := users.UserInfo{
		Username:     req.Username,
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
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	user, err := users.GetUser(ctx, req.Username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := users.VerifyPassword(req.Password, user.PasswordHash); err != nil {
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
