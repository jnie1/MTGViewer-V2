package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/users"
)

func signup(c *gin.Context) {
	var request users.SignupRequest

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := users.GetUser(request.Email); err != sql.ErrNoRows {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	passwordHash, err := users.GenerateHash(request.Password)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newUser := users.UserInfo{
		Name:         request.Name,
		Email:        request.Email,
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

	user, err := users.GetUser(request.Email)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := users.VerifyPassword(request.Password, user.PasswordHash); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	loginDuration := time.Now().Add(2 * time.Hour)

	userClaims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: loginDuration.Unix(),
			Subject:   user.Email,
		},
		Role: user.Role,
	}

	token, err := auth.GenerateToken(&userClaims)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie("token", token, int(loginDuration.Unix()), "", "", false, true)
	c.Status(http.StatusNoContent)
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.Status(http.StatusNoContent)
}

func AddUserRoutes(router *gin.Engine) {
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.POST("/logout", logout)
}
