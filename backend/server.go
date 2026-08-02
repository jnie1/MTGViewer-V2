package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/jnie1/MTGViewer-V2/routes"
)

func RegisterRouter() {
	//Opening database connection
	if err := database.Open(); err != nil {
		log.Fatal("Error opening database: ", err)
	}

	//Close database connection when RegisterRouter() returns
	//Should only happen on shutdown
	defer database.Close()

	r := gin.Default()
	r.Use(auth.CorsMiddleware())

	api := r.Group("/api")
	routes.AddUserRoutes(api)

	//Applied authorization for every route registered under this
	protected := api.Group("", auth.IsAuthorized)
	routes.AddCardRoutes(protected)
	routes.AddContainerRoutes(protected)
	routes.AddTransactionRoutes(protected)

	if err := routes.AddStaticRoutes(r, "/api"); err != nil {
		log.Fatal("Error adding static files: ", err)
	}

	r.Run(":8080")
}
