package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/jnie1/MTGViewer-V2/routes"
)

func registerRouter() {
	db, err := database.Open()
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer db.Close()

	sdk, err := cards.OpenSDK()
	if err != nil {
		log.Fatal("Error open mtg json sdk: ", err)
	}

	defer sdk.Close()

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

	go syncCards()
	if err := r.Run(":8080"); err != nil {
		log.Println("Error running: ", err)
	}
}

func syncCards() {
	ctx := context.Background()

	if err := cards.RefreshViews(ctx); err != nil {
		log.Println("issue refreshing views: ", err)
	}

	for range time.Tick(time.Hour * 24) {
		if err := cards.RefreshViews(ctx); err != nil {
			log.Println("issue refreshing views: ", err)
		}
	}
}
