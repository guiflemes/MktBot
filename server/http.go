package server

import (
	"log"
	"marketingBot/fb/handlers"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RunHttpServer(addr string) {
	app := fiber.New()
	setMiddlewares(app)

	fbApp := handlers.NewFBHttpApp()

	apiV1 := app.Group("api/v1/facebook")

	apiV1.Post("/webhook", fbApp.HandleWebhook)
	apiV1.Get("/webhook", fbApp.HandleVerification)
	apiV1.Post("/upload", fbApp.HandleUploadImage)

	log.Println("Starting HTTP server", addr)
	app.Listen(addr)

}

func setMiddlewares(app *fiber.App) {
	addCorsMiddleware(app)
	addLoggingMiddleware(app)
}

func addCorsMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
}

func addLoggingMiddleware(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Next:         nil,
		Done:         nil,
		Format:       "[${time}] ${latency} | ${path} ${status} - ${method} \n",
		TimeFormat:   "02-Jan-2006",
		TimeZone:     "America/Sao_Paulo",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	}))
}
