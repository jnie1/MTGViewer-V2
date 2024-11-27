package routes

import "github.com/gin-gonic/gin"

func signup(c *gin.Context) {
}

func login(c *gin.Context) {
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
}

func AddUserRoutes(router *gin.Engine) {
}
