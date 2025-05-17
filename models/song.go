package models

type SongPoint struct {
	SongID      string  `bson:"song_id"` // YouTube ID of the song
	Fingerprint int64   `bson:"fingerprint"`
	TimeMS      float64 `bson:"time_ms"`
}
