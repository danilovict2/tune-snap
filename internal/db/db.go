package db

import (
	"context"

	"github.com/danilovict2/shazam-clone/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SongExists(songs *mongo.Collection, songID string) bool {
	err := songs.FindOne(context.TODO(), bson.D{{Key: "song_id", Value: songID}}).Decode(&bson.M{})
	return err == nil
}

func FindSongPoints(songs *mongo.Collection, fingerprints []int64) (map[int64][]models.SongPoint, error) {
	ret := make(map[int64][]models.SongPoint)
	for _, fp := range fingerprints {
		cursor, err := songs.Find(context.TODO(), bson.D{{Key: "fingerprint", Value: fp}})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			return nil, err
		}

		var results []models.SongPoint
		if err := cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}

		ret[fp] = results
	}

	return ret, nil
}
