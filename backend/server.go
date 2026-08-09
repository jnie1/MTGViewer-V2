package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/jnie1/MTGViewer-V2/routes"
)

func RegisterRouter() {
	db, err := database.Open()
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer db.Close()

	r := gin.Default()
	r.Use(auth.CorsMiddleware())

	api := r.Group("/api")
	routes.AddUserRoutes(api)
	routes.AddCardRoutes(api)
	routes.AddContainerRoutes(api)
	routes.AddTransactionRoutes(api)

	if err := routes.AddStaticRoutes(r, "/api"); err != nil {
		log.Fatal("Error adding static files: ", err)
	}

	r.Run(":8080")
}
