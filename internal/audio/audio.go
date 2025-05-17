package audio

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/danilovict2/shazam-clone/internal/db"
	"github.com/danilovict2/shazam-clone/internal/fingerprint"
	"github.com/danilovict2/shazam-clone/internal/spotify"
	"github.com/kkdai/youtube/v2"
	"github.com/raitonoberu/ytsearch"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type searchResult struct {
	videoID    string
	durationMS int64
}

const maxRetryAttempts = 5

func SaveTracks(tracks []spotify.Track, songs *mongo.Collection) (saved int) {
	saved = len(tracks)
	var wg sync.WaitGroup
	errChan := make(chan error, len(tracks))

	go func() {
		for err := range errChan {
			log.Println(err)
			saved--
		}
	}()

	for _, track := range tracks {
		id := db.CreateSongID(track.Name, track.Artists)
		if db.SongExists(songs, id) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := downloadTrack(track, id); err != nil {
				errChan <- err
				return
			}

			wavFile, err := os.Open(filepath.Join(os.Getenv("SONGS_DIR"), id + ".wav"))
			if err != nil {
				errChan <- err
				return
			}

			song, err := ReadWav(wavFile)
			if err != nil {
				errChan <- err
				return
			}

			points, err := fingerprint.Fingerprint(song.Audio, song.Duration, song.SampleRate, id)
			if err != nil {
				errChan <- err
				return
			}

			if _, err := songs.InsertMany(context.TODO(), points); err != nil {
				errChan <- err
				return
			}
		}()
	}

	wg.Wait()
	close(errChan)
	return saved
}

func downloadTrack(track spotify.Track, outputPath string) error {
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

	ext := filepath.Ext(outputPath)
	fName := strings.TrimRight(outputPath, ext) + ".m4a"
	file, err := os.CreateTemp(os.Getenv("SONGS_DIR"), fName)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())

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

	return convertToWav(file.Name())
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
