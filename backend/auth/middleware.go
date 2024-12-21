package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AddCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", os.Getenv("CLIENT_ORIGIN"))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Next()
}

func IsAuthorized(c *gin.Context) {
	cookie, err := c.Cookie("token")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := ParseToken(cookie)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := token.Valid(); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("role", token.Role)
	c.Next()
}
