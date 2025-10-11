package handlers

import (
	"encoding/json"
	"fmt"
	"go_backend/config"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TODO: Save MBID
type RecommendationWithArtist struct {
	RecName    string `json:"recName"`
	ArtistName string `json:"artistName"`
}

type RecommendationResponse struct {
	Type    string                                `json:"type"`    // "track", "album", or "artist"
	Query   string                                `json:"query"`   // the tag or search term
	Results map[string][]RecommendationWithArtist `json:"results"` // can be []RecommendationWithArtist or map[string][]RecommendationWithArtist
	Error   string                                `json:"error,omitempty"`
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

type TopTagsResp struct {
	TopTags struct {
		Tag []struct {
			Name string `json:"name"`
		} `json:"tag"`
	} `json:"toptags"`
}

var httpClient = &http.Client{}
var baseURL = "http://ws.audioscrobbler.com/2.0/"

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
	tags := []string{}
	if response.TrackName != nil {
		params := url.Values{}
		params.Add("method", "track.getTopTags")
		params.Add("artist", strings.ToLower(response.ArtistName))
		params.Add("track", strings.ToLower(*response.TrackName))
		params.Add("api_key", config.LastFMAPIKey)
		params.Add("format", "json")
		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(
			"GET",
			fullURL,
			// "http://ws.audioscrobbler.com/2.0/"+
			// 	"?method="+"track.getTopTags"+
			// 	"&artist="+strings.ToLower(response.ArtistName)+
			// 	"&track="+strings.ToLower(*response.TrackName)+
			// 	"&api_key="+config.LastFMAPIKey+
			// 	"&format=json",
			nil,
		)
		fmt.Print("request: ", req)
		if err != nil {
			return nil, fmt.Errorf("error making http request for getting lastfm tags for track: %w", err)
		}
		var topTags TopTagsResp
		err = httpReqHelper(req, &topTags)
		if err != nil {
			return nil, fmt.Errorf("error getting lastfm tags for track: %w", err)
		}

		for _, tag := range topTags.TopTags.Tag {
			tags = append(tags, tag.Name)
		}
	} else if response.AlbumName != nil {
		params := url.Values{}
		params.Add("method", "album.getTopTags")
		params.Add("artist", strings.ToLower(response.ArtistName))
		params.Add("album", strings.ToLower(*response.AlbumName))
		params.Add("api_key", config.LastFMAPIKey)
		params.Add("format", "json")
		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(
			"GET",
			fullURL,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making http request for getting lastfm tags for track: %w", err)
		}
		var topTags TopTagsResp
		err = httpReqHelper(req, &topTags)
		if err != nil {
			return nil, fmt.Errorf("error getting lastfm tags for album: %w", err)
		}

		for _, tag := range topTags.TopTags.Tag {
			tags = append(tags, tag.Name)
		}
	} else {
		params := url.Values{}
		params.Add("method", "artist.getTopTags")
		params.Add("artist", strings.ToLower(response.ArtistName))
		params.Add("api_key", config.LastFMAPIKey)
		params.Add("format", "json")
		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(
			"GET",
			fullURL,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error making http request for getting lastfm tags for track: %w", err)
		}
		var topTags TopTagsResp
		err = httpReqHelper(req, &topTags)
		if err != nil {
			return nil, fmt.Errorf("error getting lastfm tags for album: %w", err)
		}

		for _, tag := range topTags.TopTags.Tag {
			tags = append(tags, tag.Name)
		}
	}

	return tags, nil
}

// given track name and artist name, return track recommendations from LastFM
// TODO: maybe dont take all of them? limit query?
func GetLastFMSimilarTracksBuiltin(trackName string, artistName string) (*RecommendationResponse, error) {
	params := url.Values{}
	params.Add("method", "track.getsimilar")
	params.Add("artist", strings.ToLower(artistName))
	params.Add("track", strings.ToLower(trackName))
	params.Add("api_key", config.LastFMAPIKey)
	params.Add("format", "json")
	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest(
		"GET",
		fullURL,
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
	trackRecs := make(map[string][]RecommendationWithArtist)
	trackRecs["similar tracks"] = trackNames

	return &RecommendationResponse{
		Type:    "track",
		Query:   "similar tracks",
		Results: trackRecs,
	}, nil
}

// given artist name, return artist recs from LastFM
// TODO: maybe dont take all of them? limit query?
func GetLastFMSimilarArtistsBuiltin(artistName string) (*RecommendationResponse, error) {
	params := url.Values{}
	params.Add("method", "artist.getsimilar")
	params.Add("artist", strings.ToLower(artistName))
	params.Add("api_key", config.LastFMAPIKey)
	params.Add("format", "json")
	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest(
		"GET",
		fullURL,
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

	artistNames := make([]RecommendationWithArtist, len(res.SimilarArtists.Artist))
	for index, artist := range res.SimilarArtists.Artist {
		artistNames[index] = RecommendationWithArtist{
			RecName: artist.Name,
		}
	}
	artistRecs := make(map[string][]RecommendationWithArtist)
	artistRecs["similar artists"] = artistNames

	return &RecommendationResponse{
		Type:    "artist",
		Query:   "similar artists",
		Results: artistRecs,
	}, nil
}

// given resp map and tags list, return updated resp map
// func GetLastFMRecsByTags(resp map[string][]string, tags []string) (map[string][]string, error) {
// 	resp["Track"] = GetLastFMTracksByTag(tags)
// 	resp["Album"] = GetLastFMAlbumsByTag(tags)
// 	resp["Artist"] = GetLastFMArtistsByTag(tags)
// 	return resp, nil
// }

// given list of 3 tags, return track recommendations from LastFM
func GetLastFMTracksByTag(tags []string) (*RecommendationResponse, error) {
	trackRecs := make(map[string][]RecommendationWithArtist)
	for _, tag := range tags {
		// TODO change all urls to this
		params := url.Values{}
		params.Add("method", "tag.gettoptracks")
		params.Add("tag", strings.ToLower(tag))
		params.Add("api_key", config.LastFMAPIKey)
		params.Add("format", "json")
		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(
			"GET",
			fullURL,
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

	return &RecommendationResponse{
		Type:    "track",
		Query:   strings.Join(tags, ", "),
		Results: trackRecs,
	}, nil
}

// given list of tags, return album recommendations from LastFM
func GetLastFMAlbumsByTag(tags []string) (*RecommendationResponse, error) {
	albumRecs := make(map[string][]RecommendationWithArtist)
	for _, tag := range tags {
		params := url.Values{}
		params.Add("method", "tag.gettopalbums")
		params.Add("tag", strings.ToLower(tag))
		params.Add("api_key", config.LastFMAPIKey)
		params.Add("format", "json")
		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(
			"GET",
			fullURL,
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

	return &RecommendationResponse{
		Type:    "album",
		Query:   strings.Join(tags, ", "),
		Results: albumRecs,
	}, nil
}

// given list of tags, return artist recommendations from LastFM
func GetLastFMArtistsByTag(tags []string) (*RecommendationResponse, error) {
	artistRecs := make(map[string][]RecommendationWithArtist)
	for _, tag := range tags {
		params := url.Values{}
		params.Add("method", "tag.gettopartists")
		params.Add("tag", strings.ToLower(tag))
		params.Add("api_key", config.LastFMAPIKey)
		params.Add("format", "json")
		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(
			"GET",
			fullURL,
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

		artistRecsForTag := make([]RecommendationWithArtist, len(res.Artists.Artist))
		for index, recommendedArtist := range res.Artists.Artist {
			artistRecsForTag[index] = RecommendationWithArtist{
				RecName:    recommendedArtist.Name,
				ArtistName: "", // no artist field for artist recs
			}
		}
		artistRecs[tag] = artistRecsForTag
	}

	return &RecommendationResponse{
		Type:    "artist",
		Query:   strings.Join(tags, ", "),
		Results: artistRecs,
	}, nil
}

func httpReqHelper(httpReq *http.Request, varForUnmarshal any) error {
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error making http client/sending req: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resp body: %w", err)
	}

	err = json.Unmarshal(body, varForUnmarshal)
	if err != nil {
		return fmt.Errorf("error unmarshalling resp body: %w", err)
	}

	return nil
}
