package cors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupCors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                            // TODO change after frontend hosted
		AllowMethods: "GET,POST,PUT,DELETE",          // Allowed methods
		AllowHeaders: "Origin, Content-Type, Accept", // Allowed headers
	}))
}
