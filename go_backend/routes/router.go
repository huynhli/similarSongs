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
	respObj, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}
	tags, err := handlers.GetLastFMTags(respObj)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM tags.",
		}
	}

	lastFmRecs, err := handlers.GetLastFMTracksByTag(tags)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM tracks by tag.",
		}
	}
	return c.JSON(lastFmRecs)
}

func GetLastFMArtist(c *fiber.Ctx) error {
	link := c.Query("link")
	respObj, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	tags, err := handlers.GetLastFMTags(respObj)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM tags.",
		}
	}

	lastFmRecs, err := handlers.GetLastFMArtistsByTag(tags)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM artists by tag.",
		}
	}
	return c.JSON(lastFmRecs)
}

func GetLastFMAlbum(c *fiber.Ctx) error {
	link := c.Query("link")
	respObj, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	tags, err := handlers.GetLastFMTags(respObj)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM tags.",
		}
	}

	lastFmRecs, err := handlers.GetLastFMAlbumsByTag(tags)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM albums by tag.",
		}
	}
	return c.JSON(lastFmRecs)
}

func GetLastFMTrackSimilar(c *fiber.Ctx) error {
	link := c.Query("link")
	respObj, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	lastFmRecs, err := handlers.GetLastFMSimilarTracksBuiltin(*respObj.TrackName, respObj.ArtistName)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM similar tracks.",
		}
	}
	return c.JSON(lastFmRecs)
}

func GetLastFMArtistSimilar(c *fiber.Ctx) error {
	link := c.Query("link")
	respObj, err := handlers.SpotifyInNameOut(link)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error with spotifyAPIHandler.",
		}
	}

	lastFmRecs, err := handlers.GetLastFMSimilarArtistsBuiltin(respObj.ArtistName)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error getting LastFM similar artists.",
		}
	}
	return c.JSON(lastFmRecs)
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
