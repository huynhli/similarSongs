package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	SpotifyClientID     string
	SpotifyClientSecret string
	LastFMAPIKey        string
	LastFMSharedSecret  string
	// SpotifyRedirectURI  string
)

func LoadConfig() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set config variables from environment variables
	SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	LastFMAPIKey = os.Getenv("LAST_FM_API_KEY")
	LastFMSharedSecret = os.Getenv("LAST_FM_SECRET")
	// SpotifyRedirectURI = "https://song-recommendations-web-app.pages.dev/auth/callback" // Redirect URI for your app TODO fix when frontend hosted
}
