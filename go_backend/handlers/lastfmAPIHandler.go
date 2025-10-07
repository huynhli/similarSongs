package handlers

import (
	"encoding/json"
	"fmt"
	"go_backend/config"
	"io"
	"net/http"
	"strings"
)

// TODO: Save MBID
type RecommendationWithArtist struct {
	RecName    string
	ArtistName string
}

type LastFMTrackSimilar struct {
	SimilarTracks struct {
		Track []struct {
			Name   string  `json:"name"`
			Match  float64 `json:"match"`
			MBID   string  `json:"mbid"`
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"track"`
	} `json:"similartracks"`
}

type LastFMArtistSimilar struct {
	SimilarArtists struct {
		Artist []struct {
			Name  string `json:"name"`
			Match string `json:"match"`
			MBID  string `json:"mbid"`
		} `json:"artist"`
	} `json:"similarartists"`
}

type LastFMTrackRec struct {
	Tracks struct {
		Track []struct {
			Name   string `json:"name"`
			MBID   string `json:"mbid"`
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"track"`
	} `json:"tracks"`
}

type LastFMArtistRec struct {
	Artists struct {
		Artist []struct {
			Name string `json:"name"`
			MBID string `json:"mbid"`
		} `json:"artist"`
	} `json:"topartists"`
}

type LastFMAlbumRec struct {
	Albums struct {
		Album []struct {
			Name   string `json:"name"`
			MBID   string `json:"mbid"`
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"album"`
	} `json:"albums"`
}

type TopTags struct {
	Tag []struct {
		Name string `json:"name"`
	}
}

var httpClient = &http.Client{}

// take in spotify response obj, return list of recommendations (artists + albums + track)
func GetLastFMTags(response SpotifyResponseObj) ([]string, error) {
	// track -- artist tags, album tags, track tags, track similar
	// album -- artist tags, album tags
	// artist -- artist tags, artists similar

	// ensure lastfm apikey

	// // map with arist, track, album keys
	// resp := make(map[string][]string)
	// // resp["Artist"] = []string{}
	// // resp["ArtistSimilar"] = []string{}
	// // resp["Track"] = []string{}
	// // resp["TrackSimilar"] = []string{}
	// // resp["Album"] = []string{}

	// TODO: get tags
	tags := []string{}
	if response.TrackName != nil {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"track.getTopTags"+
				"&artist="+strings.ToLower(response.ArtistName)+
				"&track="+strings.ToLower(*response.TrackName)+
				"&api_key="+config.LastFMAPIKey+
				"&format=json",
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making http request for getting lastfm tags for track: %w", err)
		}
		var topTags TopTags
		err = httpReqHelper(req, &topTags)
		if err != nil {
			return nil, fmt.Errorf("error getting lastfm tags for track: %w", err)
		}

		for _, tag := range topTags.Tag {
			tags = append(tags, tag.Name)
		}
	} else if response.AlbumName != nil {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"album.getTopTags"+
				"&artist="+strings.ToLower(response.ArtistName)+
				"&track="+strings.ToLower(*response.AlbumName)+
				"&api_key="+config.LastFMAPIKey+
				"&format=json",
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making http request for getting lastfm tags for track: %w", err)
		}
		var topTags TopTags
		err = httpReqHelper(req, &topTags)
		if err != nil {
			return nil, fmt.Errorf("error getting lastfm tags for album: %w", err)
		}

		for _, tag := range topTags.Tag {
			tags = append(tags, tag.Name)
		}
	} else {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"artist.getTopTags"+
				"&artist="+strings.ToLower(response.ArtistName)+
				"&api_key="+config.LastFMAPIKey+
				"&format=json",
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making http request for getting lastfm tags for track: %w", err)
		}
		var topTags TopTags
		err = httpReqHelper(req, &topTags)
		if err != nil {
			return nil, fmt.Errorf("error getting lastfm tags for album: %w", err)
		}

		for _, tag := range topTags.Tag {
			tags = append(tags, tag.Name)
		}
	}

	return tags, nil
}

// given track name and artist name, return track recommendations from LastFM
// TODO: maybe dont take all of them? limit query?
func GetLastFMSimilarTracksBuiltin(trackName string, artistName string) ([]RecommendationWithArtist, error) {
	req, err := http.NewRequest(
		"GET",
		"http://ws.audioscrobbler.com/2.0/"+
			"?method="+"track.getsimilar"+
			"&track="+strings.ToLower(trackName)+
			"&artist="+strings.ToLower(artistName)+
			"&api_key="+config.LastFMAPIKey+
			"&format=json",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error making similar tracks http request: %w", err)
	}

	var res LastFMTrackSimilar
	err = httpReqHelper(req, &res)
	if err != nil {
		return nil, err
	}

	trackNames := make([]RecommendationWithArtist, len(res.SimilarTracks.Track))
	for index, track := range res.SimilarTracks.Track {
		trackNames[index] = RecommendationWithArtist{
			RecName:    track.Name,
			ArtistName: track.Artist.Name,
		}
	}

	return trackNames, nil
}

// given artist name, return artist recs from LastFM
// TODO: maybe dont take all of them? limit query?
func GetLastFMSimilarArtistsBuiltin(artistName string) ([]string, error) {
	req, err := http.NewRequest(
		"GET",
		"http://ws.audioscrobbler.com/2.0/"+
			"?method="+"artist.getsimilar"+
			"&artist="+strings.ToLower(artistName)+
			"&api_key="+config.LastFMAPIKey+
			"&format=json",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error making similar artist http request: %w", err)
	}

	var res LastFMArtistSimilar
	err = httpReqHelper(req, &res)
	if err != nil {
		return nil, err
	}

	artistNames := make([]string, len(res.SimilarArtists.Artist))
	for index, artist := range res.SimilarArtists.Artist {
		artistNames[index] = artist.Name
	}

	return artistNames, nil
}

// given resp map and tags list, return updated resp map
// func GetLastFMRecsByTags(resp map[string][]string, tags []string) (map[string][]string, error) {
// 	resp["Track"] = GetLastFMTracksByTag(tags)
// 	resp["Album"] = GetLastFMAlbumsByTag(tags)
// 	resp["Artist"] = GetLastFMArtistsByTag(tags)
// 	return resp, nil
// }

// given list of 3 tags, return track recommendations from LastFM
func GetLastFMTracksByTag(tags []string) (map[string][]RecommendationWithArtist, error) {
	trackRecs := make(map[string][]RecommendationWithArtist)
	for _, tag := range tags {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"tag.gettoptracks"+
				"&tag="+strings.ToLower(tag)+
				"&api_key="+config.LastFMAPIKey+
				"&format=json",
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making track recs by tag http request: %w", err)
		}

		var res LastFMTrackRec
		err = httpReqHelper(req, &res)
		if err != nil {
			return nil, err
		}

		trackRecsForTag := make([]RecommendationWithArtist, len(res.Tracks.Track))
		for index, recommendedTrack := range res.Tracks.Track {
			trackRecsForTag[index] = RecommendationWithArtist{
				RecName:    recommendedTrack.Name,
				ArtistName: recommendedTrack.Artist.Name,
			}
		}
		trackRecs[tag] = trackRecsForTag
	}

	return trackRecs, nil
}

// given list of tags, return album recommendations from LastFM
func GetLastFMAlbumsByTag(tags []string) (map[string][]RecommendationWithArtist, error) {
	albumRecs := make(map[string][]RecommendationWithArtist)
	for _, tag := range tags {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"tag.gettopalbums"+
				"&tag="+strings.ToLower(tag)+
				"&api_key="+config.LastFMAPIKey+
				"&format=json",
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making album recs by tag http request: %w", err)
		}

		var res LastFMAlbumRec
		err = httpReqHelper(req, &res)
		if err != nil {
			return nil, err
		}

		albumRecsForTag := make([]RecommendationWithArtist, len(res.Albums.Album))
		for index, recommendedAlbum := range res.Albums.Album {
			albumRecsForTag[index] = RecommendationWithArtist{
				RecName:    recommendedAlbum.Name,
				ArtistName: recommendedAlbum.Artist.Name,
			}
		}
		albumRecs[tag] = albumRecsForTag
	}

	return albumRecs, nil
}

// given list of tags, return artist recommendations from LastFM
func GetLastFMArtistsByTag(tags []string) (map[string][]string, error) {
	artistRecs := make(map[string][]string)
	for _, tag := range tags {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"tag.gettopartists"+
				"&tag="+strings.ToLower(tag)+
				"&api_key="+config.LastFMAPIKey+
				"&format=json",
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making album recs by tag http request: %w", err)
		}

		var res LastFMArtistRec
		err = httpReqHelper(req, &res)
		if err != nil {
			return nil, err
		}

		artistRecsForTag := make([]string, len(res.Artists.Artist))
		for index, recommendedArtist := range res.Artists.Artist {
			artistRecsForTag[index] = recommendedArtist.Name
		}
		artistRecs[tag] = artistRecsForTag
	}

	return artistRecs, nil
}

func httpReqHelper(httpReq *http.Request, varForUnmarshal any) error {
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error making similar artist http client/sending req: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resp body of similar artist: %w", err)
	}

	err = json.Unmarshal(body, varForUnmarshal)
	if err != nil {
		return fmt.Errorf("error unmarshalling resp body of similar artist: %w", err)
	}

	return nil
}
