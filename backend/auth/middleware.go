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
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions},
		AllowHeaders:     []string{"Content-type", "Authorization"},
		AllowCredentials: true,
	})
}

func IsAuthorized(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if !strings.HasPrefix(token, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := ParseToken(token)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := ValidateClaims(claims); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("role", claims.Role)
	c.Next()
}
