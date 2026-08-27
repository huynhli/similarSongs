package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
)

func LoadConfig() {
	// Load environment variables from .env file
	if os.Getenv("DOCKER") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, relying on system environment variables")
		}
	}
	// Set config variables from environment variables
	SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	// SpotifyRedirectURI = "https://song-recommendations-web-app.pages.dev/auth/callback"
	SpotifyRedirectURI = "http://localhost:5173/auth/callback"
}
