package shazam

import (
	"log"
	"math"
	"sort"
	"sync"

	"github.com/danilovict2/shazam-clone/internal/db"
	"github.com/danilovict2/shazam-clone/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Match struct {
	SongID string
	Score  int
}

func Recognize(recordedPoints []models.SongPoint, songs *mongo.Collection) ([]Match, error) {
	log.Printf("Recognize: received %d recorded points", len(recordedPoints))
	fingerprintMap := toFingeprintMap(recordedPoints)
	dbPoints := getDbPoints(fingerprintMap, songs)
	log.Printf("Recognize: fetched song points for %d fingerprints from DB", len(dbPoints))

	timestamps := make(map[string][][2]float64, 0)

	for fp, points := range dbPoints {
		for _, point := range points {
			timestamps[point.SongID] = append(timestamps[point.SongID], [2]float64{point.TimeMS, fingerprintMap[fp].TimeMS})
		}
	}
	log.Printf("Recognize: built timestamps map for %d songs", len(timestamps))

	// SongID -> Score
	matchesMap := make(map[string]int)
	for songID, songTimestamps := range timestamps {
		score := 0
		for i := range songTimestamps {
			for j := i + 1; j < len(songTimestamps); j++ {
				dbDiff := math.Abs(songTimestamps[i][0] - songTimestamps[j][0])
				recordedDiff := math.Abs(songTimestamps[i][1] - songTimestamps[j][1])

				// Allow a small tolerance
				if math.Abs(dbDiff-recordedDiff) < 100 {
					score++
				}
			}
		}
		matchesMap[songID] = score
		log.Printf("Recognize: songID=%s, score=%d", songID, score)
	}

	matches := make([]Match, 0, len(matchesMap))
	for songID, score := range matchesMap {
		matches = append(matches, Match{SongID: songID, Score: score})
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	log.Printf("Recognize: returning %d matches", len(matches))
	return matches, nil
}

func toFingeprintMap(points []models.SongPoint) map[int64]models.SongPoint {
	ret := make(map[int64]models.SongPoint)
	for _, point := range points {
		ret[point.Fingerprint] = point
	}

	return ret
}

func getDbPoints(fingerprints map[int64]models.SongPoint, songs *mongo.Collection) map[int64][]models.SongPoint {
	dbPoints := make(map[int64][]models.SongPoint)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for fp := range fingerprints {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			points, err := db.SongPointsWithFingerprint(fp, songs)
			if err != nil {
				log.Printf("getDbPoints: failed to fetch song points for fingerprint %d: %v", fp, err)
				return
			}

			mu.Lock()
			dbPoints[fp] = points
			mu.Unlock()
		}()
	}

	wg.Wait()
	return dbPoints
}
