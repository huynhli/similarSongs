package handlers

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetRecommendationsAPI(c *fiber.Ctx) error {
	rand.Seed(time.Now().UnixNano())

	// if an error occured
	returnedGenres, token := GetGenreFromAPI(c)
	if len(returnedGenres) == 0 {
		temp := []string{"This artist/track/album has not been categorized by Spotify. Or, the playlist you entered does not have any tracks."}
		return c.JSON(temp)
	} else if strings.HasPrefix(returnedGenres[0], "This") {
		return c.JSON(returnedGenres)
	}
	// else, get recommendations!
	// TODO this is deprecated -> find a new way to do this

	totalSongsMap := make(map[string]struct{})
	totalSongsList := []string{}

	supportedGenresList := []string{}
	playlistIDs := []string{}
	for _, eachGenre := range returnedGenres {
		// TODO throw playlist links in a file and call dynamically
		switch eachGenre {
		case "rap":
			playlistIDs = append(playlistIDs, "7k8E4tmsRVxLSQitB9bYPQ")
			playlistIDs = append(playlistIDs, "01MRi9jFGeSEEttKOk7VgR")
			playlistIDs = append(playlistIDs, "23bBJex1i75b9oqUo6jUCC")
			playlistIDs = append(playlistIDs, "7JIGfa0KkCTDxUPOQySODP")
			supportedGenresList = append(supportedGenresList, "rap")
		case "hip hop":
			playlistIDs = append(playlistIDs, "62y3BHKehWnb1hlaPclDAA")
			playlistIDs = append(playlistIDs, "5oN7X3cPTUkOJFPmVx5wCE")
			playlistIDs = append(playlistIDs, "0dMexqq0XIWS3QJ74z3ZhD")
			playlistIDs = append(playlistIDs, "6WnTjZBYrYOdPl4pK7PyZg")
			supportedGenresList = append(supportedGenresList, "hip hop")
		case "pop":
			playlistIDs = append(playlistIDs, "1WH6WVBwPBz35ZbWsgCpgr")
			playlistIDs = append(playlistIDs, "3JoHkM90TXzfIS1RMN0Cgd")
			playlistIDs = append(playlistIDs, "2L2HwKRvUgBv1YetudaRI3")
			playlistIDs = append(playlistIDs, "2Sxd5BovYwLgRg6KyZOTer")
			supportedGenresList = append(supportedGenresList, "pop")
		case "rock":
			playlistIDs = append(playlistIDs, "1Kgkkup7w9qvGxGJGa75PS")
			playlistIDs = append(playlistIDs, "77RvyLiqmUimojxq3vg6mY")
			playlistIDs = append(playlistIDs, "23hD5D7bvXtkJGz2ni7s9e")
			playlistIDs = append(playlistIDs, "7KahejEXeb27cRJCFt3VFO")
			supportedGenresList = append(supportedGenresList, "rock")
		}
	}

	for _, eachPlaylistId := range playlistIDs {
		// getting offset num
		reqNum, errNum := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+eachPlaylistId, nil)
		if errNum != nil {
			log.Fatalf("Error getting playlist: %v", errNum)
		}
		reqNum.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		respNum, errNum := client.Do(reqNum)
		if errNum != nil {
			log.Fatalf("Error sending request: %v", errNum)
		}
		defer respNum.Body.Close()

		bodyNum, errNum := io.ReadAll(respNum.Body)
		if errNum != nil {
			log.Fatalf("Error reading response body: %v", errNum)
		}
		var Result struct {
			Tracks struct {
				Total int `json:"total"`
			} `json:"tracks"`
		}

		errNum = json.Unmarshal(bodyNum, &Result)
		if errNum != nil {
			log.Fatalf("Error unmarshalling: %v", errNum)
		}

		// calculate offset for this playlist
		offset := rand.Intn(Result.Tracks.Total - 5)
		if Result.Tracks.Total <= 5 {
			offset = 0
		}

		// fetching 5 tracks from each playlist with offset
		req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+eachPlaylistId+"/tracks?offset="+strconv.Itoa(offset)+"&limit=5", nil)
		if err != nil {
			log.Fatalf("Error getting recommendations: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer respNum.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		var Playlist struct {
			Items []struct {
				TrackObject struct {
					Artists []struct {
						Name string `json:"name"`
					} `json:"artists"`
					Name string `json:"name"`
				} `json:"track"`
				Name string `json:"name"`
			} `json:"items"`
		}

		err = json.Unmarshal(body, &Playlist)
		if err != nil {
			log.Fatalf("Error unmarshalling: %v", err)
		}

		for _, eachItem := range Playlist.Items {
			var ArtistList = ""
			for indexArtist, eachArtist := range eachItem.TrackObject.Artists {
				if indexArtist == 0 {
					ArtistList = ArtistList + " by " + eachArtist.Name
					continue
				}
				ArtistList = ArtistList + " + " + eachArtist.Name
			}
			addUnique(totalSongsMap, eachItem.TrackObject.Name+" "+ArtistList+" ", &totalSongsList)
		}
	}
	// unique songs in the list (no duplicates)
	totalSongsList = append(totalSongsList, strings.Join(returnedGenres, ", ")+".", strings.Join(supportedGenresList, ", ")+".")

	// randomly pick 15
	songsToDisplay := 16
	if len(totalSongsList) <= songsToDisplay+2 {
		return c.JSON(totalSongsList)
	} else {
		// fisher yates shuffle (i just learned this)

		n := len(totalSongsList) - 2
		for i := 0; i < n-1; i++ {
			j := rand.Intn(n-i) + i
			totalSongsList[i], totalSongsList[j] = totalSongsList[j], totalSongsList[i] // Swap
		}
		listToSend := make([]string, songsToDisplay)
		copy(listToSend, totalSongsList[:songsToDisplay])
		listToSend = append(listToSend, totalSongsList[n:]...)
		return c.JSON(listToSend)
	}
}
