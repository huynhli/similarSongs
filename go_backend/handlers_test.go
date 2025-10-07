package main

import (
	"go_backend/config"
	"go_backend/handlers"
	"go_backend/routes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetLastFMTrackRoute(t *testing.T) {
	config.LoadConfig()
	app := fiber.New()
	app.Get("/test", routes.GetLastFMTrack)

	req := httptest.NewRequest("GET", "/test?link=https://open.spotify.com/artist/6vWDO969PvNqNYHIOW5v0m?si=i8pS1ntxTF-tCTV1rItksw", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)

	// Check that the JSON contains the expected fields
	assert.Contains(t, string(body), `"artist"`)
	// assert.Contains(t, string(body), `"track"`)
	// assert.Contains(t, string(body), `"album"`)
}

func TestSpotifyAndLastFMHandlers(t *testing.T) {
	// TODO group tests
}

// TODO mockable tests

func TestSpotifyInNameOutArtist(t *testing.T) {
	config.LoadConfig()
	resp, err := handlers.SpotifyInNameOut("https://open.spotify.com/artist/6vWDO969PvNqNYHIOW5v0m?si=i8pS1ntxTF-tCTV1rItksw")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ArtistName)
	assert.Equal(t, "Beyoncé", resp.ArtistName)
}

func TestSpotifyInNameOutTrack(t *testing.T) {
	config.LoadConfig()
	resp, err := handlers.SpotifyInNameOut("https://open.spotify.com/track/0TwBtDAWpkpM3srywFVOV5")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ArtistName)
	assert.NotEmpty(t, resp.TrackName)
	assert.NotEmpty(t, resp.AlbumName)
	assert.Equal(t, "Beyoncé", resp.ArtistName)
	// assert.Equal(t, []string{"Beyoncé", "JAY-Z"}, resp.ArtistName)
	assert.Equal(t, "Crazy In Love (feat. JAY-Z)", *resp.TrackName)
	assert.Equal(t, "Dangerously In Love", *resp.AlbumName)
}

func TestSpotifyInNameOutAlbum(t *testing.T) {
	config.LoadConfig()
	resp, err := handlers.SpotifyInNameOut("https://open.spotify.com/album/75sBENwOba7T6buiOJsAZs?si=3xw4VPFAQ3mbXHq3EqFtfQ")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ArtistName)
	assert.NotEmpty(t, resp.AlbumName)
	assert.Equal(t, "Long Live The Empire", resp.ArtistName)
	assert.Equal(t, "Deathless", *resp.AlbumName)
}

func TestLastFMSimilarTracksBuiltin(t *testing.T) {
	config.LoadConfig()
	resp, err := handlers.GetLastFMSimilarTracksBuiltin("believe", "cher")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, "Cher", resp[0].ArtistName)
	assert.Equal(t, "Strong Enough", resp[0].RecName)
}

func TestLastFMSimilarArtistsBuiltin(t *testing.T) {
	config.LoadConfig()
	resp, err := handlers.GetLastFMSimilarArtistsBuiltin("cher")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, "Sonny & Cher", resp[0])
}

func TestLastFMTracksByTag(t *testing.T) {
	config.LoadConfig()
	tag := []string{"disco"}
	resp, err := handlers.GetLastFMTracksByTag(tag)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, "Dancing Queen", resp[tag[0]][0].RecName)
	assert.Equal(t, "ABBA", resp[tag[0]][0].ArtistName)
}

func TestLastFMAlbumByTag(t *testing.T) {
	config.LoadConfig()
	tag := []string{"disco"}
	resp, err := handlers.GetLastFMAlbumsByTag(tag)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, "Bad Girls", resp[tag[0]][0].RecName)
	assert.Equal(t, "Donna Summer", resp[tag[0]][0].ArtistName)
}

func TestLastFMArtistByTag(t *testing.T) {
	config.LoadConfig()
	tag := []string{"disco"}
	resp, err := handlers.GetLastFMArtistsByTag(tag)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, "ABBA", resp[tag[0]][0])
}

// func TestSpotifyInNameOutAlbum(t *testing.T) {
// 	config.LoadConfig()
// 	app := fiber.New()
// 	app.Get("/test", handlers.InToOut)

// 	req := httptest.NewRequest("GET", "/test?link=https://open.spotify.com/album/6BzxX6zkDsYKFJ04ziU5xQ?si=cchby_lyTES2MnA1L2kleg", nil)
// 	resp, err := app.Test(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)

// 	body, _ := io.ReadAll(resp.Body)

// 	// Check that the JSON contains the expected fields
// 	assert.Contains(t, string(body), `"artist"`)
// 	// assert.Contains(t, string(body), `"track"`)
// 	assert.Contains(t, string(body), `"album"`)
// }

// func TestSpotifyInNameOutTrack(t *testing.T) {
// 	config.LoadConfig()
// 	app := fiber.New()
// 	app.Get("/test", handlers.InToOut)

// 	req := httptest.NewRequest("GET", "/test?link=https://open.spotify.com/track/0TwBtDAWpkpM3srywFVOV5?si=2c1aa537c2fe4340", nil)
// 	resp, err := app.Test(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)

// 	body, _ := io.ReadAll(resp.Body)

// 	// Check that the JSON contains the expected fields
// 	assert.Contains(t, string(body), `"artist"`)
// 	assert.Contains(t, string(body), `"track"`)
// 	assert.Contains(t, string(body), `"album"`)
// }

// TODO more tests
