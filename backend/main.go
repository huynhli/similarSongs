package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"go_backend/config"
	"go_backend/cors"
	"go_backend/routes"
)

func init() {
	config.LoadConfig()
}

func main() {

	//Set port from env var
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}

	//create new Fiber app
	app := fiber.New()

	//enable CORS to allow requests to backend
	cors.SetupCors(app)

	//routing
	routes.SetupRoutes(app)

	//port listen and serve
	err := app.Listen(":" + port)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
