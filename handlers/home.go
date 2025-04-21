package handlers

import (
	"net/http"

	"github.com/danilovict2/shazam-clone/templates/home"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	return Render(c, http.StatusOK, home.Hello())
}