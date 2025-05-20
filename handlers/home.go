package handlers

import (
	"log"
	"net/http"

	"github.com/danilovict2/shazam-clone/internal/db"
	"github.com/danilovict2/shazam-clone/templates/home"
	"github.com/labstack/echo/v4"
)

func (cfg *Config) Home(c echo.Context) error {
	count, err := db.GetSongCount(cfg.MongoClient.Database("shazam").Collection("songs"))
	if err != nil {
		log.Print(err)
		return err
	}

	return Render(c, http.StatusOK, home.Hello(count))
}
