package youtube

import (
	"fmt"
	"strings"

	"github.com/danilovict2/shazam-clone/internal/spotify"
	"github.com/raitonoberu/ytsearch"
)

type searchResult struct {
	videoID    string
	durationMS int64
}

func DownloadTracks(tracks []spotify.Track) error {
	for _, track := range tracks {
		if err := downloadTrack(track); err != nil {
			return err
		}
	}

	return nil
}

func downloadTrack(track spotify.Track) error {
	_, err := findBestMatch(track)
	return err
}

func findBestMatch(track spotify.Track) (searchResult, error) {
	results, err := search(track)
	if err != nil {
		return searchResult{}, err
	}

	fmt.Println(results)
	return searchResult{}, nil
}

func search(track spotify.Track) ([]searchResult, error) {
	searchQuery := fmt.Sprintf("'%s' %s", track.Name, strings.Join(track.Artists, " "))
	search := ytsearch.VideoSearch(searchQuery)

	result, err := search.Next()
	if err != nil {
		return nil, err
	}

	results := make([]searchResult, 0)
	for i := range min(len(result.Videos), 10) {
		video := result.Videos[i]
		if video != nil {
			results = append(results, searchResult{videoID: video.ID, durationMS: int64(video.Duration * 1000)})
		}
	}

	return results, nil
}
