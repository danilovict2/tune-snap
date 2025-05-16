package shazam

import (
	"math"

	"github.com/danilovict2/shazam-clone/internal/db"
	"github.com/danilovict2/shazam-clone/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Recognize(recordedPoints []models.SongPoint, songs *mongo.Collection) (map[string]int, error) {
	fingerprintMap := toFingeprintMap(recordedPoints)

	fingerprints := make([]int64, 0)
	for fp := range fingerprintMap {
		if fp != 0 {
			fingerprints = append(fingerprints, fp)
		}
	}

	dbPoints, err := db.FindSongPoints(songs, fingerprints)
	if err != nil {
		return nil, err
	}

	timestamps := make(map[string][][2]float64, 0)

	for fp, points := range dbPoints {
		for _, point := range points {
			timestamps[point.SongID] = append(timestamps[point.SongID], [2]float64{point.TimeMS, fingerprintMap[fp].TimeMS})
		}
	}

	// SongID -> Score
	matches := make(map[string]int)
	for songID, songTimestamps := range timestamps {
		for i := range songTimestamps {
			for j := i + 1; j < len(songTimestamps); j++ {
				dbDiff := math.Abs(songTimestamps[i][0] - songTimestamps[j][0])
				recordedDiff := math.Abs(songTimestamps[i][1] - songTimestamps[j][1])

				// Allow a small tolerance
				if math.Abs(dbDiff - recordedDiff) < 100 {
					matches[songID]++
				}
			}
		}
	}

	return matches, nil
}

func toFingeprintMap(points []models.SongPoint) map[int64]models.SongPoint {
	ret := make(map[int64]models.SongPoint)
	for _, point := range points {
		ret[point.Fingerprint] = point
	}

	return ret
}
