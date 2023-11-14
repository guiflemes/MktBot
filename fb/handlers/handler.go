package handlers

import (
	"log"
	"marketingBot/fb/adapters"
	"marketingBot/fb/models"
	"marketingBot/fb/service"
	"marketingBot/settings"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type (
	Authenticator interface {
		Auth(c *fiber.Ctx) error
	}

	MessageHandler interface {
		HandleWebHookRequest(models.WehbookReq) error
	}

	ImageRequest struct {
		Url string `json:"url"`
	}
)

type FBHttpApp struct {
	auth       Authenticator
	msgHandler MessageHandler
	uplodImage func(url string) (string, error)
}

func NewFBHttpApp() *FBHttpApp {
	return &FBHttpApp{
		auth:       NewPageAcesssAuth(),
		msgHandler: service.NewSimpleMessageUC(),
		uplodImage: adapters.UploadImage,
	}
}

func (fb *FBHttpApp) HandleVerification(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == settings.GETENV("FB_VERIFY_TOKEN") {
		return c.Status(http.StatusOK).SendString(challenge)
	}

	return c.SendStatus(http.StatusForbidden)

}

func (fb *FBHttpApp) HandleWebhook(c *fiber.Ctx) error {
	// err := fb.auth.Auth(c)

	// if err != nil {
	// 	log.Println("unauthorized", err)
	// 	return c.Status(http.StatusUnauthorized).SendString("unauthorized")
	// }

	var webhookReq models.WehbookReq

	err := c.BodyParser(&webhookReq)

	if err != nil {
		log.Println("body parser request", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to body parser request"})
	}

	return fb.msgHandler.HandleWebHookRequest(webhookReq)
}

func (fb *FBHttpApp) HandleUploadImage(c *fiber.Ctx) error {

	var imageReq ImageRequest

	if err := c.BodyParser(imageReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to body parser request"})
	}

	attachmentID, err := fb.uplodImage(imageReq.Url)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload image"})
	}

	return c.JSON(fiber.Map{"attachment_id": attachmentID})
}
