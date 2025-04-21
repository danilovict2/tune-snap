package main

import (
	"log"
	"os"

	"github.com/danilovict2/shazam-clone/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()

	e.GET("/", handlers.Home)

	e.Logger.Fatal(e.Start(os.Getenv("LISTEN_ADDR")))
}
