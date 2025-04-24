package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/shazam-clone/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

	e.GET("/", handlers.Home)
	
	api := e.Group("/api")
	api.POST("/recognize", handlers.Recognize)

	e.Logger.Fatal(e.Start(os.Getenv("LISTEN_ADDR")))
}
