package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	var allowedOrigins []string
	originsEnv := os.Getenv("CLIENT_ORIGINS")

	for origin := range strings.SplitSeq(originsEnv, ",") {
		trimmedOrigin := strings.TrimSpace(origin)
		if trimmedOrigin != "" {
			allowedOrigins = append(allowedOrigins, trimmedOrigin)
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-type", "Authorization"},
		AllowCredentials: true,
	})
}

func IsAuthorized(c *gin.Context) {
	creds := c.GetHeader("Authorization")

	if !strings.HasPrefix(creds, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	creds = strings.TrimPrefix(creds, "Bearer ")
	token, err := ParseToken(creds)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := ValidateClaims(token); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("role", token.Role)
	c.Next()
}
