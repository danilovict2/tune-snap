package spotify

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	trackEndpoint    = "https://api.spotify.com/v1/tracks/"
	albumEndpoint    = "https://api.spotify.com/v1/albums/"
	playlistEndpoint = "https://api.spotify.com/v1/playlists/"
)

func GetTracks(url, token string) ([]Track, error) {
	// Remove the trailing slash if the user included it in the URL
	url = strings.TrimRight(url, "/")
	switch {
	case strings.Contains(url, "track"):
		return extractSingleTrack(url, token)
	case strings.Contains(url, "album"):
		return extractAlbumTracks(url, token)
	case strings.Contains(url, "playlist"):
		return extractPlaylistTracks(url, token)
	default:
		return nil, fmt.Errorf("invalid url: %s", url)
	}
}

func extractSingleTrack(url, token string) ([]Track, error) {
	trackPattern := `^https:\/\/open\.spotify\.com\/track\/[a-zA-Z0-9]{22}$`
	re := regexp.MustCompile(trackPattern)
	if !re.MatchString(url) {
		return nil, fmt.Errorf("invalid track url: %s", url)
	}

	id := extractIDFromUrl(url)
	json, err := apiRequest(trackEndpoint+id, token)
	if err != nil {
		return nil, err
	}

	return []Track{trackInfo(json)}, nil
}

func extractAlbumTracks(url, token string) ([]Track, error) {
	albumPattern := `^https:\/\/open\.spotify\.com\/album\/[a-zA-Z0-9]{22}$`
	re := regexp.MustCompile(albumPattern)
	if !re.MatchString(url) {
		return nil, fmt.Errorf("invalid album url: %s", url)
	}

	id := extractIDFromUrl(url)
	json, err := apiRequest(albumEndpoint+id, token)
	if err != nil {
		return nil, err
	}

	tracks := make([]Track, 0)
	result := gjson.Get(json, "tracks.items")
	result.ForEach(func(key, value gjson.Result) bool {
		tracks = append(tracks, trackInfo(value.String()))
		return true
	})

	return tracks, nil
}

func extractPlaylistTracks(url, token string) ([]Track, error) {
	playlistPattern := `^https:\/\/open\.spotify\.com\/playlist\/[a-zA-Z0-9]{22}$`
	re := regexp.MustCompile(playlistPattern)
	if !re.MatchString(url) {
		return nil, fmt.Errorf("invalid album url: %s", url)
	}

	id := extractIDFromUrl(url)
	json, err := apiRequest(playlistEndpoint+id, token)
	if err != nil {
		return nil, err
	}

	tracks := make([]Track, 0)
	result := gjson.Get(json, "tracks.items.#.track")
	result.ForEach(func(key, value gjson.Result) bool {
		tracks = append(tracks, trackInfo(value.String()))
		return true
	})

	return tracks, nil
}

func extractIDFromUrl(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}

func apiRequest(url, token string) (respBody string, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", resp.Status)
	}

	json, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

func trackInfo(json string) Track {
	track := Track{Name: gjson.Get(json, "name").String(), DurationMS: gjson.Get(json, "duration_ms").Int()}
	artists := gjson.Get(json, "artists.#.name")
	for _, artist := range artists.Array() {
		track.Artists = append(track.Artists, artist.String())
	}

	return track
}
