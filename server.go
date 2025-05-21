package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/tune-snap/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connection to database: %v", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	config := handlers.Config{MongoClient: client}

	e := echo.New()
	e.Static("/public", "public")
	
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))


	e.GET("/", config.Home)
	
	api := e.Group("/api")
	api.POST("/recognize", config.Recognize)
	api.POST("/add_song", config.AddSong)
	api.GET("/spotify_auth", config.SpotifyAuth)


	e.Logger.Fatal(e.Start(os.Getenv("LISTEN_ADDR")))
}