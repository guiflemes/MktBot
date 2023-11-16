package server

import (
	"fmt"
	"log"
	fb "marketingBot/fb/handlers"
	"os"
	"time"

	dash "marketingBot/dashboard/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RunHttpServer(port string) {

	if port == "" {
		log.Fatal("port must be set")
	}

	app := fiber.New()
	setMiddlewares(app)

	fbApp := fb.NewFBHttpApp()
	dashApp := dash.NewdashHttpApp()
	flowApp := dash.NewFlowApp()

	apiV1 := app.Group("api/v1")

	apiV1.Post("/facebook/webhook", fbApp.HandleWebhook)
	apiV1.Get("/facebook/webhook", fbApp.HandleVerification)
	apiV1.Post("/facebook/upload", fbApp.HandleUploadImage)

	apiV1.Get("/dash/clicks", dashApp.HandleClickCount)
	apiV1.Get("/dash/revels", dashApp.HandleCouponRevelCount)

	apiV1.Get("/flow/:key", flowApp.HandleGetFLow)
	apiV1.Post("/flow", flowApp.HandleSaveFlow)

	log.Println("Starting HTTP server", port)
	app.Listen(fmt.Sprintf(":%s", port))

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
