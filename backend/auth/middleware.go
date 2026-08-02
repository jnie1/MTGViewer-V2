package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	originsEnv := os.Getenv("CLIENT_ORIGINS")
	allowedOrigins := []string{}

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

// Used to restrict routes to users whose JWT role is "admin"
// Assumes IsAuthorized() is already called
func RequireAdmin(c *gin.Context) {
	role, exists := c.Get("role")

	//exists == false, means that IsAuthorized() didn't run or never set a role
	//role != "admin", means a valid user logged in but they just aren't an admin
	//http.StatusForbidden is given, because the user is authenticated, they just lack permission
	if !exists || role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}
