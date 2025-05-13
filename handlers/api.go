package handlers

import (
	"io"
	"net/http"
	"strconv"

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

	_, err = io.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = strconv.ParseFloat(c.FormValue("audio_duration"), 64)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
