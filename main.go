package main

import (
	"marketingBot/server"
	"marketingBot/settings"
)

func main() {
	server.RunHttpServer(settings.GETENV("PORT"))
}
