package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"go_backend/config"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Artist struct {
	Name string `json:"name"`
}

type Album struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Track struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
	Album   Album    `json:"album"`
}

type SpotifyResponseObj struct {
	ArtistName string  `json:"artist"`
	AlbumName  *string `json:"album,omitempty"`
	TrackName  *string `json:"track,omitempty"`
}

// input a artist, album, or track. get a json obj with Artist, ?Album, ?Track.
func SpotifyInNameOut(link string) (SpotifyResponseObj, error) {
	// get trimmed, lowercase, valid link as input, along with type of output wanted
	// TODO any link validation to still do (edge cases)
	// build request with sanitizing/validating

	linkSplit := strings.Split(link, "/")
	if len(linkSplit) < 5 {
        return nil, fmt.Errorf("invalid link: %s", link)
    }
	typeLink := linkSplit[3]

	var decider string
	switch typeLink {
	case "artist":
		decider = "artists"
	case "album":
		decider = "albums"
	case "track":
		decider = "tracks"
	default:
		return nil, fmt.Errorf("Link does not point to an artist, album, or track. Please try a different link, and not: %s", link)
	}

	// valid link, so get token
	spotifyAccessToken, err := getSpotifyAPIToken()
	if err != nil {
		return nil, fmt.Errorf("Token error: %s", err)
	}

	// finish building req
	// var requestURL strings.Builder
	// artist --> artists/
	// album --> albums/
	// track --> tracks/
	// make proper request, return obj of artist name, ?album name, ?track name
	id := strings.Split(linkSplit[4], "?")
	requestURL := `https://api.spotify.com/v1/` + decider + "/" + id[0]
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error calling API: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+spotifyAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request to Spotify API: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body from Spotify API: %s", err)
	}

	log.Printf("Spotify raw response: %s", string(body))
	var response SpotifyResponseObj
	switch decider {
	case "artists":
		var a Artist
		err := json.Unmarshal(body, &a)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling artist from Spotify API: %s", err)
		}
		response.ArtistName = a.Name

	case "albums":
		var a Album
		err := json.Unmarshal(body, &a)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling album from Spotify API: %s", err)
		}
		response.ArtistName = a.Artists[0].Name
		response.AlbumName = &a.Name

	case "tracks":
		var t Artist
		err := json.Unmarshal(body, &t)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling track from Spotify API: %s", err)
		}
		response.ArtistName = t.Artists[0].Name
		response.AlbumName = &t.Album.Name
		response.TrackName = &t.Name
	}

	// artist --> returns artist
	// album --> returns album + artist
	// track --> returns track + album + artist
	return response, nil
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// generates token to use with spotify api
func getSpotifyAPIToken() (string, error) {
	// pull secret
	clientID := config.SpotifyClientID
	clientSecret := config.SpotifyClientSecret

	// build request
	authString := clientID + ":" + clientSecret
	authStringBase64 := base64.StdEncoding.EncodeToString([]byte(authString))

	url := "https://accounts.spotify.com/api/token"
	data := "grant_type=client_credentials"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "Error building Spotify's POST request", err
	}

	req.Header.Set("Authorization", "Basic "+authStringBase64)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// init http client + send req
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Error sending POST request to Spotify", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading POST response body form Spotify", err
	}

	// unmarshal into struct/var
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "Error unmarshalling POST JSON response into Token", err
	}

	return tokenResponse.AccessToken, nil
}

// returns popularity and genre of spotify artist
func spotifyArtistPopAndGenre(spotifyArtistIDs []string) (int, []string, error) {
	genres := []string{"pop", "dance", spotifyArtistIDs[0]}
	popularity := 100
	return popularity, genres, nil
}
