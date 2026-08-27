package routes

import (
	"go_backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Register routes for authentication
	app.Get("/", homePage)
	app.Get("/api/data", handlers.GetRecommendationsAPI) //TODO change to GetRecommendationsAPI
	app.Get("/api/cmsMaintenance", handlers.GetMaintenanceStatus)
}

func homePage(c *fiber.Ctx) error {
	return c.SendString("Hi this is the home page of the song recommendation web app. The Github repo can be found at: https://github.com/huynhli/song-recommendations-web-app")
}

// func getGenreAPI(c *fiber.Ctx) error {
// 	link := c.Query("link")
// 	tempList := []string{"This is a valid link.", "This is not a valid link. Try again."}
// 	if link == "" {
// 		return c.JSON(tempList[1:])
// 	}

// 	// authOptions = []string{}

// 	return c.JSON(tempList[:1])
// }
