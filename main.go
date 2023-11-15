package main

import (
	"log"
	"marketingBot/server"
	"marketingBot/settings"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server.RunHttpServer(settings.GETENV("PORT"))
}
