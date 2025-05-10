package youtube

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/danilovict2/shazam-clone/internal/spotify"
	"github.com/kkdai/youtube/v2"
	"github.com/raitonoberu/ytsearch"
)

type searchResult struct {
	videoID    string
	durationMS int64
}

const maxRetryAttempts = 5

func DownloadTracks(tracks []spotify.Track) error {
	for _, track := range tracks {
		if err := downloadTrack(track); err != nil {
			return err
		}
	}

	return nil
}

func downloadTrack(track spotify.Track) error {
	match, err := findBestMatch(track)
	if err != nil {
		return err
	}

	client := youtube.Client{}

	video, err := client.GetVideo(match.videoID)
	if err != nil {
		return err
	}

	formats := video.Formats.Itag(140)

	fName := fmt.Sprintf("%s_%s", track.Name, strings.Join(track.Artists, "_")) + ".m4a"
	file, err := os.Create(os.Getenv("SONGS_DIR") + "/" + fName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Downloads may occasionally fail when using the github.com/kkdai/youtube/v2 library.
	// Implement a retry mechanism with a limit on the number of attempts to handle such cases.
	var isDownloaded bool = false
	attempt := 1

	for !isDownloaded {
		if attempt > maxRetryAttempts {
			return fmt.Errorf("failed to download video: %s", video.Title)
		}

		stream, _, err := client.GetStream(video, &formats[0])
		if err != nil {
			return err
		}
		defer stream.Close()

		if _, err := io.Copy(file, stream); err != nil {
			return err
		}

		fi, err := file.Stat()
		if err != nil {
			return err
		}

		attempt++
		isDownloaded = fi.Size() > 0
	}

	return nil
}

func findBestMatch(track spotify.Track) (searchResult, error) {
	results, err := search(track)
	if err != nil {
		return searchResult{}, err
	}

	var (
		minDiff float64 = 1e9
		ret     searchResult
	)

	for _, result := range results {
		if diff := math.Abs(float64(result.durationMS - track.DurationMS)); diff < minDiff {
			ret = result
			minDiff = diff
		}
	}

	return ret, nil
}

func search(track spotify.Track) ([]searchResult, error) {
	searchQuery := fmt.Sprintf("'%s' %s", track.Name, strings.Join(track.Artists, " "))
	search := ytsearch.VideoSearch(searchQuery)

	result, err := search.Next()
	if err != nil {
		return nil, err
	}

	if len(result.Videos) == 0 {
		return nil, fmt.Errorf("no videos found for query: %s", searchQuery)
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
