package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Config struct {
	MongoClient        *mongo.Client
	SpotifyAccessToken string
}

func Render(c echo.Context, code int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(code, buf.String())
}
