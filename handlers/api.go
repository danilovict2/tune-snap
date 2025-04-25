package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/danilovict2/shazam-clone/internal/fingerprint"
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

	fmt.Println(fingerprint.Fingerprint(sample, audioDuration))

	return c.String(http.StatusOK, "")
}
