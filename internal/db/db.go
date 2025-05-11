package db

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SongExists(songs *mongo.Collection, songID string) bool {
	err := songs.FindOne(context.TODO(), bson.D{{Key: "SongID", Value: songID}}).Decode(&bson.M{})
	return err == nil
}

func CreateSongID(name string, artists []string) string {
	id := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	for _, artist := range artists {
		id += "_" + strings.ToLower(strings.ReplaceAll(artist, " ", "_"))
	}

	return id
}
