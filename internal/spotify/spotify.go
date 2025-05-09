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
	trackEndpoint = "https://api.spotify.com/v1/tracks/"
	albumEndpoint = "https://api.spotify.com/v1/albums/"
)

func GetTracks(url, token string) ([]Track, error) {
	// Remove the trailing slash if the user included it in the URL
	url = strings.TrimRight(url, "/")
	switch {
	case strings.Contains(url, "track"):
		return extractSingleTrack(url, token)
	case strings.Contains(url, "album"):
		return extractAlbumTracks(url, token)
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
	req, err := http.NewRequest(http.MethodGet, trackEndpoint+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	json, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	track := Track{Name: gjson.Get(string(json), "name").String()}
	artists := gjson.Get(string(json), "artists.#.name")
	for _, artist := range artists.Array() {
		track.Artists = append(track.Artists, artist.String())
	}

	return []Track{track}, nil
}

func extractAlbumTracks(url, token string) ([]Track, error) {
	albumPattern := `^https:\/\/open\.spotify\.com\/album\/[a-zA-Z0-9]{22}$`
	re := regexp.MustCompile(albumPattern)
	if !re.MatchString(url) {
		return nil, fmt.Errorf("invalid album url: %s", url)
	}

	id := extractIDFromUrl(url)
	req, err := http.NewRequest(http.MethodGet, albumEndpoint+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	json, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tracks := make([]Track, 0)
	result := gjson.Get(string(json), "tracks.items")
	result.ForEach(func(key, value gjson.Result) bool {
		track := Track{Name: gjson.Get(value.String(), "name").String()}
		artists := gjson.Get(value.String(), "artists.#.name")
		for _, artist := range artists.Array() {
			track.Artists = append(track.Artists, artist.String())
		}

		tracks = append(tracks, track)
		return true
	})

	fmt.Println(tracks)
	return tracks, nil
}

func extractIDFromUrl(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}
