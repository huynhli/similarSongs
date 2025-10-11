package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/huynhli/similarsongs/go_backend/config"
	"github.com/huynhli/similarsongs/go_backend/handlers"
	"github.com/huynhli/similarsongs/go_backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetLastFMTrackRoute(t *testing.T) {
	config.LoadConfig()
	app := fiber.New()
	app.Get("/test", routes.GetLastFMTrack)

	req := httptest.NewRequest("GET", "/test?link=https://open.spotify.com/artist/3l0CmX0FuQjFxr8SK7Vqag?si=JZ_FrvKmTJ6waWoL8J_peQ", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"artistName"`)

	var rec handlers.RecommendationResponse
	_ = json.Unmarshal(body, &rec)

	assert.Equal(t, "Space Song", rec.Results["indie pop"][0].RecName)
	assert.Equal(t, "Beach House", rec.Results["indie pop"][0].ArtistName)
	assert.Equal(t, "Sofia", rec.Results["indie pop"][1].RecName)
	assert.Equal(t, "Clairo", rec.Results["indie pop"][1].ArtistName)
	assert.Equal(t, "Buzzcut Season", rec.Results["indie pop"][2].RecName)
	assert.Equal(t, "Lorde", rec.Results["indie pop"][2].ArtistName)

	assert.Equal(t, "Bags", rec.Results["bedroom pop"][0].RecName)
	assert.Equal(t, "Clairo", rec.Results["bedroom pop"][0].ArtistName)
	assert.Equal(t, "Pretty Girl", rec.Results["bedroom pop"][1].RecName)
	assert.Equal(t, "Clairo", rec.Results["bedroom pop"][1].ArtistName)
	assert.Equal(t, "Bad Habit", rec.Results["bedroom pop"][2].RecName)
	assert.Equal(t, "Steve Lacy", rec.Results["bedroom pop"][2].ArtistName)

	assert.Equal(t, "Apocalypse", rec.Results["dream pop"][0].RecName)
	assert.Equal(t, "Cigarettes After Sex", rec.Results["dream pop"][0].ArtistName)
	assert.Equal(t, "Space Song", rec.Results["dream pop"][1].RecName)
	assert.Equal(t, "Beach House", rec.Results["dream pop"][1].ArtistName)
	assert.Equal(t, "The Subway", rec.Results["dream pop"][2].RecName)
	assert.Equal(t, "Chappell Roan", rec.Results["dream pop"][2].ArtistName)

	// Check that the JSON contains the expected fields
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
	// assert.NotEmpty(t, resp.AlbumName)
	assert.Equal(t, "Beyoncé", resp.ArtistName)
	// assert.Equal(t, []string{"Beyoncé", "JAY-Z"}, resp.ArtistName)
	assert.Equal(t, "Crazy In Love (feat. JAY-Z)", *resp.TrackName)
	// assert.Equal(t, "Dangerously In Love", *resp.AlbumName)
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

	assert.Equal(t, "Strong Enough", resp.Results["similar tracks"][0].RecName)
	assert.Equal(t, "Cher", resp.Results["similar tracks"][0].ArtistName)

}

func TestLastFMSimilarArtistsBuiltin(t *testing.T) {
	config.LoadConfig()
	resp, err := handlers.GetLastFMSimilarArtistsBuiltin("cher")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	assert.Equal(t, "Sonny & Cher", resp.Results["similar artists"][0].RecName)
}

func TestLastFMTracksByTag(t *testing.T) {
	config.LoadConfig()
	tag := []string{"disco"}
	resp, err := handlers.GetLastFMTracksByTag(tag)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	assert.Equal(t, "September", resp.Results[tag[0]][0].RecName)
	assert.Equal(t, "Earth, Wind & Fire", resp.Results[tag[0]][0].ArtistName)
}

func TestLastFMAlbumByTag(t *testing.T) {
	config.LoadConfig()
	tag := []string{"disco"}
	resp, err := handlers.GetLastFMAlbumsByTag(tag)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	assert.Equal(t, "Off the Wall", resp.Results[tag[0]][0].RecName)
	assert.Equal(t, "Michael Jackson", resp.Results[tag[0]][0].ArtistName)
}

func TestLastFMArtistByTag(t *testing.T) {
	config.LoadConfig()
	tag := []string{"disco"}
	resp, err := handlers.GetLastFMArtistsByTag(tag)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	assert.Equal(t, "ABBA", resp.Results[tag[0]][0].RecName)
}

func TestLastFMTrackTags(t *testing.T) {
	config.LoadConfig()
	trackTemp := "Paranoid Android"

	testResp := handlers.SpotifyResponseObj{
		ArtistName: "Radiohead",
		TrackName:  &trackTemp,
	}
	tags, err := handlers.GetLastFMTags(testResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tags)
	assert.Equal(t, "alternative rock", tags[0])
	assert.Equal(t, "alternative", tags[1])
	assert.Equal(t, "rock", tags[2])
}

func TestLastFMAlbumTags(t *testing.T) {
	config.LoadConfig()
	albumTemp := "The Bends"

	testResp := handlers.SpotifyResponseObj{
		ArtistName: "Radiohead",
		AlbumName:  &albumTemp,
	}
	tags, err := handlers.GetLastFMTags(testResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tags)
	assert.Equal(t, "alternative rock", tags[0])
	assert.Equal(t, "1995", tags[1])
	assert.Equal(t, "rock", tags[2])
}

func TestLastFMArtistTags(t *testing.T) {
	config.LoadConfig()
	testResp := handlers.SpotifyResponseObj{
		ArtistName: "Clairo",
	}
	tags, err := handlers.GetLastFMTags(testResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tags)
	assert.Equal(t, "indie pop", tags[0])
	assert.Equal(t, "bedroom pop", tags[1])
	assert.Equal(t, "dream pop", tags[2])
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
