package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/jnie1/MTGViewer-V2/routes"
)

func RegisterRouter() {
	if err := database.Open(); err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer database.Close()

	r := gin.Default()
	r.Use(auth.CorsMiddleware())

	if err := routes.AddStaticPaths(r); err != nil {
		log.Fatal("Error adding static files: ", err)
	}

	api := r.Group("/api")
	routes.AddUserRoutes(api)
	routes.AddCardRoutes(api)
	routes.AddContainerRoutes(api)
	routes.AddTransactionRoutes(api)

	authorized := api.Group("", auth.IsAuthorized)

	authorized.GET("/secret", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"secret": "some secret"})
	})

	r.Run(":8080")
}
