package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedOrigins []string

// Function runs one time upon initialization, before main() gets called
func init() {
	origins := os.Getenv("CLIENT_ORIGINS")

	//Parse the CLIENT_ORIGINS string into multiple tokens and append it to the list of valid Origins
	for _, origin := range strings.Split(origins, ",") {
		trimmedOrigin := strings.TrimSpace(origin)
		if trimmedOrigin != "" {
			allowedOrigins = append(allowedOrigins, trimmedOrigin)
		}
	}
}

func AddCors(c *gin.Context) {
	requestOrigin := c.Request.Header.Get("Origin")

	for _, origin := range allowedOrigins {
		if origin == requestOrigin {
			c.Header("Access-Control-Allow-Origin", requestOrigin)
			break
		}
	}

	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
	c.Header("Vary", "Origin")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

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
