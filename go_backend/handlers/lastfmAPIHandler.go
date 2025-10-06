package handlers

import (
	"encoding/json"
	"fmt"
	"go_backend/config"
)

// take in spotify response obj, return list of recommendations (artists + albums + track)
func GetLastFMRecsFromSpotify(response SpotifyResponseObj) ([]byte, error) {
	// track -- artist tags, album tags, track tags, track similar
	// album -- artist tags, album tags
	// artist -- artist tags, artists similar

	// ensure lastfm apikey


	// map with arist, track, album keys
	resp := make(map[string][]string)
	// resp["Artist"] = []string{}
	// resp["ArtistSimilar"] = []string{}
	// resp["Track"] = []string{}
	// resp["TrackSimilar"] = []string{}
	// resp["Album"] = []string{}

	if response.TrackName != nil {
		// track similar
		// get similar tracks using track name, artist name
		resp["TrackSimilar"] = GetLastFMSimilarTracksBuiltin(*response.TrackName, response.ArtistName) 
	}
	else if response.TrackName != nil {
		// track similar
		// get similar tracks using track name, artist name
		resp["AlbumSimilar"] = GetLastFMSimilarTracksBuiltin(*response.TrackName, response.ArtistName) 
	}
	else if response.AlbumName != nil {
		// artists similar
		resp["ArtistSimilar"] = GetLastFMSimilarArtistsBuiltin(response.ArtistName)
	}

	// TODO: get tags
	tags := []string{}

	// use tags by case, so probably diff function/endpoint. maybe based on response?
	resp = GetLastFMRecsByTags(resp, tags)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshalling LastFM response recs: %w", err)
	}

	return jsonResp, nil
}

// given track name and artist name, return track recommendations from LastFM
// TODO: maybe dont take all of them? limit query?
type LastFMTrackSimilar struct {
	SimilarTracks struct {
		Track []struct {
			Name 	string	`json:"name"`
			Match	float64	`json:"match"`
			MBID	string	`json:"mbid"`
			Artist 	struct {
				Name	string	`json:"name"`	
			} `json:"artist"`
		} `json:"track"`
	} `json:"similartracks"`
}

func GetLastFMSimilarTracksBuiltin(trackName string, artistName string) ([]string, error) {
	req, err := http.NewRequest(
		"GET",
		"http://ws.audioscrobbler.com/2.0/"+
			"?method="+"track.getsimilar"
			+"&track="+trackName
			+"&artist="+artistName
			+"&api_key="+config.LastFMAPIKey
			+"&format=json", 
		nil
	)
	if err != nil {
		return nil, fmt.Errorf("error making similar tracks http request: %w", err)
	}

	resp, err := &http.Client{}.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making similar tracks http client/sending req: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading resp body of similar tracks: %w", err)
	}

	var res LastFMTrackSimilar 
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling resp body of similar tracks: %w", err)
	}

	trackNames := make([](string, string), len(res.SimilarTracks.Track))
	for index, track := range res.SimilarTracks.Track {
		trackNames[index] = (track.Name, track.Artist.Name)
	}

	return trackNames, nil
}

// given artist name, return artist recs from LastFM
// TODO: maybe dont take all of them? limit query?
type LastFMArtistSimilar struct {
	SimilarArtists struct {
		Track []struct {
			Name 	string	`json:"name"`
			Match	float64	`json:"match"`
			MBID	string	`json:"mbid"`
		} `json:"artist"`
	} `json:"similarartists"`
}

func GetLastFMSimilarArtistsBuiltin(artistName string) ([]string, error) {
	req, err := http.NewRequest(
		"GET",
		"http://ws.audioscrobbler.com/2.0/"+
			"?method="+"artist.getsimilar"+
			"&artist="+artistName
			+"&api_key="+config.LastFMAPIKey
			+"&format=json", 
		nil
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

func httpReqHelper(httpReq *http.Request, varForUnmarshal any) error {
	resp, err := &http.Client{}.Do(httpReq)
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

// given resp map and tags list, return updated resp map
func GetLastFMRecsByTags(resp map[string][]string, tags []string) (map[string][]string, error) {
	resp["Track"] = GetLastFMTracksByTag(tags)
	resp["Album"] = GetLastFMAlbumsByTag(tags)
	resp["Artist"] = GetLastFMArtistsByTag(tags)
	return resp, nil
}

// given list of 3 tags, return track recommendations from LastFM
type LastFMTrackRec struct {
	Tracks struct {
		Track []struct {
			Name 	string	`json:"name"`
			MBID	string	`json:"mbid"`
		} `json:"track"`
	} `json:"tracks"`
}

func GetLastFMTracksByTag(tags []string) ([]string, error) {
	trackRecs := make(map[string][]string)
	for _, tag := range tags {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"tag.gettoptracks"+
				"&tag="+tag
				+"&api_key="+config.LastFMAPIKey
				+"&format=json",
			nil
		)
		if err != nil {
			return nil, fmt.Errorf("error making track recs by tag http request: %w", err)
		}

		resp, err := &http.Client{}.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making track recs by tag http client/sending req: %w", err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading resp body of track recs by tag: %w", err)
		}

		var res LastFMTrackRec
		err := json.Unmarshal(body, &res)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling resp body of track recs by tag: %w", err)
		}

		trackRecsForTag := make([]string, len(res.Tracks.Track))
		for index, recommendedTrack := range res.Tracks.Track {
			trackRecsForTag[index] = recommendedTrack
		}
		trackRecs[tag] = trackRecsForTag
	}

	return trackRecs, nil
}

// given list of tags, return album recommendations from LastFM
type LastFMAlbumRec struct {
	Albums struct {
		Album []struct {
			Name 	string	`json:"name"`
			MBID	string	`json:"mbid"`
			Artist 	struct {
				Name	string	`json:"name"`	
			} `json:"artist"`
		} `json:"album"`
	} `json:"albums"`
}

func GetLastFMAlbumsByTag(tags []string) ([]string, error) {
	albumRecs := make(map[string][]string)
	for _, tag := range tags {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"tag.gettopalbums"+
				"&tag="+tag
				+"&api_key="+config.LastFMAPIKey
				+"&format=json",
			nil
		)
		if err != nil {
			return nil, fmt.Errorf("error making album recs by tag http request: %w", err)
		}

		resp, err := &http.Client{}.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making album recs by tag http client/sending req: %w", err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading resp body of album recs by tag: %w", err)
		}

		var res LastFMAlbumRec
		err := json.Unmarshal(body, &res)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling resp body of album recs by tag: %w", err)
		}

		albumRecsForTag := make([]string, len(res.Albums.Album))
		for index, recommendedAlbum := range res.Albums.Album {
			albumRecsForTag[index] = recommendedAlbum
		}
		albumRecs[tag] = albumRecsForTag
	}

	return albumRecs, nil
}

// given list of tags, return artist recommendations from LastFM
type LastFMArtistRec struct {
	Artists struct {
		Artist []struct {
			Name 	string	`json:"name"`
			MBID	string	`json:"mbid"`
		} `json:"artist"`
	} `json:"topartists"`
}

func GetLastFMArtistsByTag(tags []string) ([]string, error) {
	artistRecs := make(map[string][]string)
	for _, tag := range tags {
		req, err := http.NewRequest(
			"GET",
			"http://ws.audioscrobbler.com/2.0/"+
				"?method="+"tag.gettopartists"+
				"&tag="+tag
				+"&api_key="+config.LastFMAPIKey
				+"&format=json",
			nil
		)
		if err != nil {
			return nil, fmt.Errorf("error making album recs by tag http request: %w", err)
		}

		resp, err := &http.Client{}.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making album recs by tag http client/sending req: %w", err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading resp body of album recs by tag: %w", err)
		}

		var res LastFMArtistRec
		err := json.Unmarshal(body, &res)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling resp body of album recs by tag: %w", err)
		}

		artistRecsForTag := make([]string, len(res.Artists.Artist))
		for index, recommendedArtist := range res.Artists.Artist {
			artistRecsForTag[index] = recommendedArtist
		}
		artistRecs[tag] = artistRecsForTag
	}

	return artistRecs, nil
}
