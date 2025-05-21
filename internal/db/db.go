package db

import (
	"context"
	"fmt"

	"github.com/danilovict2/tune-snap/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SongExists(songs *mongo.Collection, songID string) bool {
	err := songs.FindOne(context.TODO(), bson.D{{Key: "song_id", Value: songID}}).Decode(&bson.M{})
	return err == nil
}

func SongPointsWithFingerprint(fingeprint int64, songs *mongo.Collection) ([]models.SongPoint, error) {
	cursor, err := songs.Find(context.TODO(), bson.D{{Key: "fingerprint", Value: fingeprint}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var results []models.SongPoint
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetSongCount(songs *mongo.Collection) (int32, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$song_id"}}}},
		bson.D{{Key: "$count", Value: "song_count"}},
	}

	cursor, err := songs.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, err
	}

	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		return 0, err
	}

	count, ok := result[0]["song_count"].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected song_count type: %T", result[0]["song_count"])
	}

	return count, nil
}
