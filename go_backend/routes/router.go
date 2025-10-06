package routes

import (
	"go_backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Register routes for authentication
	app.Get("/", homePage)

	api := app.Group("/api/v1")

	lastFMAPI := api.Group("/lastfm")
	lastFMAPI.Get("/track", GetLastFMTrack)
	lastFMAPI.Get("/tracksimilar", GetLastFMTrackSimilar)
	lastFMAPI.Get("/artist", GetLastFMArtist)
	lastFMAPI.Get("/artistsimilar", GetLastFMArtistSimilar)
	lastFMAPI.Get("/album", GetLastFMAlbum)

	// explore = api.Group("/explore")

	// musicBrainzAPI = explore.Group("/musicBrainz")
	// musicBrainzAPI.Get("")

	// deezerAPI := explore.Group("/deezer")
	// deezerAPI.Get("")
}

func GetLastFMTrack(c *fiber.Ctx) error {
	link := c.Query("link")
	jsonResponse, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	return c.JSON(jsonResponse)
}

func GetLastFMArtist(c *fiber.Ctx) error {
	link := c.Query("link")
	jsonResponse, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	return c.JSON(jsonResponse)
}

func GetLastFMAlbum(c *fiber.Ctx) error {
	link := c.Query("link")
	jsonResponse, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	return c.JSON(jsonResponse)
}

func GetLastFMTrackSimilar(c *fiber.Ctx) error {
	link := c.Query("link")
	jsonResponse, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	return c.JSON(jsonResponse)
}

func GetLastFMArtistSimilar(c *fiber.Ctx) error {
	link := c.Query("link")
	jsonResponse, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	return c.JSON(jsonResponse)
}

func homePage(c *fiber.Ctx) error {

	return c.SendString("Hi this is the home page of a song recommendation web app. The Github repo can be found at: https://github.com/huynhli/similarSongs")
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
