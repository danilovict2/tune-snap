package handlers

import (
	"io"
	"net/http"

	"github.com/danilovict2/shazam-clone/internal/fingerprint"
	"github.com/labstack/echo/v4"
)

func Recognize(c echo.Context) error {
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

	fingerprint.Fingerprint(sample)
	
	return c.String(http.StatusOK, "")
}
