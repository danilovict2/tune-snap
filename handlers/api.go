package handlers

import (
	"net/http"

	"github.com/danilovict2/tune-snap/internal/audio"
	"github.com/danilovict2/tune-snap/internal/fingerprint"
	"github.com/danilovict2/tune-snap/internal/shazam"
	"github.com/labstack/echo/v4"
)

func (cfg *Config) Recognize(c echo.Context) error {
	header, err := c.FormFile("sample")
	if err != nil {
		return err
	}

	f, err := header.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	sample, err := audio.ReadWav(f)
	if err != nil {
		return err
	}

	fingeprints, err := fingerprint.Fingerprint(sample.Audio, sample.Duration, sample.SampleRate, "")
	if err != nil {
		return err
	}

	matches, err := shazam.Recognize(fingeprints, cfg.MongoClient.Database("shazam").Collection("songs"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, matches)
}
