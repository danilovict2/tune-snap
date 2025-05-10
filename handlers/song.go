package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/danilovict2/shazam-clone/internal/spotify"
	"github.com/kashifkhan0771/utils/rand"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

const (
	spotifyUserAuthorizationEndpoint = "https://accounts.spotify.com/authorize?"
	spotifyTokenEndpoint             = "https://accounts.spotify.com/api/token"
)

func (cfg *Config) AddSong(c echo.Context) error {
	tracks, err := spotify.GetTracks(c.FormValue("url"), cfg.SpotifyAccessToken)
	if err != nil {
		// Handle 400 Bad Request error by initiating user re-authentication
		if err.Error() == fmt.Sprintf("%d %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)) {
			state, err := rand.StringWithLength(16)
			if err != nil {
				return err
			}

			query := url.Values{
				"response_type": {"code"},
				"client_id":     {os.Getenv("SPOTIFY_CLIENT_ID")},
				"scope":         {"playlist-read-private playlist-read-collaborative"},
				"redirect_uri":  {os.Getenv("SPOTIFY_REDIRECT_URI")},
				"state":         {state},
			}

			return c.Redirect(http.StatusSeeOther, spotifyUserAuthorizationEndpoint+query.Encode())
		}

		return err
	}

	return c.JSON(http.StatusOK, tracks)
}

func (cfg *Config) SpotifyAuth(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if state == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET")))

	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", os.Getenv("HOME_URI"))
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", spotifyTokenEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cfg.SpotifyAccessToken = gjson.Get(string(body), "access_token").Str

	return c.Redirect(http.StatusSeeOther, "/")
}
