package main

import (
	"log"
	"marketingBot/server"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server.RunHttpServer(os.Getenv("PORT"))
}
