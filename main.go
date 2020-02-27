package main

import (
	"log"
	"os"

	"github.com/elsagg/luzia/app"
	"github.com/joho/godotenv"
)

var appEnv = os.Getenv("APP_ENVIRONMENT")

func main() {
	if appEnv != "production" {
		_ = godotenv.Load()
		log.Println("loading .env environment")
	}

	port := os.Getenv("APP_PORT")

	server := app.NewServer()

	server.Start(port)
}
