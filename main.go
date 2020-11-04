package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"Prodkt/config"
	"Prodkt/routes"
)

func main()  {
	// Database
	config.Connect()

	// Init Router
	router := gin.Default()

	// Route Handlers / Endpoints
	routes.Routes(router)

	log.Fatal(router.Run(":8080"))
}