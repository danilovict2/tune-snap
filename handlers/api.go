package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/danilovict2/shazam-clone/internal/audio"
	"github.com/danilovict2/shazam-clone/internal/fingerprint"
	"github.com/danilovict2/shazam-clone/internal/shazam"
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

	sample, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	audioDuration, err := strconv.ParseFloat(c.FormValue("audio_duration"), 64)
	if err != nil {
		return err
	}

	samples, err := audio.BytesToSamples(sample)
	if err != nil {
		return err
	}

	sampleRate, err := strconv.ParseUint(c.FormValue("sample_rate"), 10, 32)
	if err != nil {
		return err
	}

	fingeprints, err := fingerprint.Fingerprint(samples, audioDuration, uint32(sampleRate), "")
	if err != nil {
		return err
	}

	matches, err := shazam.Recognize(fingeprints, cfg.MongoClient.Database("shazam").Collection("songs"))
	if err != nil {
		return err
	}

	fmt.Println(matches)
	return c.JSON(http.StatusOK, matches)
}
