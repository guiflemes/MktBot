package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"marketingBot/fb/models"
	"net/http"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

const (
	verifyToken     = "TES123COCO"
	pageAccessToken = "EAAMUNWFZBapoBOyZAov9RwyXzQ7DJYZC8NzvAPE2sXQuFOB905OFQVYBIozNMZAekBxZBvoaddfiI519upALZAX3oCCQ1ZBiNnYERZA1RJTYvCcgPSvOLBbpuH1qQjPFBZBHXuf5R275sO5RBYPLqdWJWVeB7vswoLxfBKmiVVdeLlfxUH2j1rD2obEhWJuZBVqnmF"
)

type Authenticator interface {
	Auth(c *fiber.Ctx) error
}

type FBHttpApp struct {
	auth Authenticator
}

func NewFBHttpApp() *FBHttpApp {
	return &FBHttpApp{
		auth: NewPageAcesssAuth(),
	}
}

func (fb *FBHttpApp) HandleVerification(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
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
		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	uc := NewSimpleMessageUC()

	return uc.HandleWebHookRequest(webhookReq)
}

func (fb *FBHttpApp) HandleUploadImage(c *fiber.Ctx) error {

	attachmentID, err := UploadImage("")

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload image"})
	}

	return c.JSON(fiber.Map{"attachment_id": attachmentID})
}

type MessageFlow struct {
	lockFlow   sync.RWMutex
	flow       map[string]models.Message
	defaultMsg models.Message
}

func NewPostbackFlow() *MessageFlow {
	flow := &MessageFlow{flow: make(map[string]models.Message)}
	flow.Add("are_you_a_dog_yes", models.Message{Attachment: &models.Attachment{
		Type: "template",
		Payload: models.PayloadCoupon{
			TemplateType: "coupon",
			Title:        "10% off everything",
			CouponCode:   "10PERCENT",
			Payload:      "coupon_10_off",
		},
	}}).
		Add("are_you_a_dog_no", models.Message{Attachment: &models.Attachment{
			Type: "template",
			Payload: models.PayloadMedia{
				TemplateType: "media",
				Elements: []models.MediaElement{
					{
						MediaType:    "image",
						AttachmentId: "342132335170214",
					},
				},
			},
		}})
	return flow
}

func NewMessageFlow() *MessageFlow {
	flow := &MessageFlow{flow: make(map[string]models.Message)}
	flow.Add("hello", models.Message{Text: "world"}).
		Add("indio", models.Message{Text: "parana"}).
		Add("rudens", models.Message{Attachment: &models.Attachment{
			Type: "template",
			Payload: models.PayloadButtons{
				TemplateType: "button",
				Text:         "Are you a dog?",
				Buttons: []models.Button{
					{
						Type:    "postback",
						Title:   "Yes",
						Payload: "are_you_a_dog_yes",
					},
					{
						Type:    "postback",
						Title:   "No",
						Payload: "are_you_a_dog_no",
					},
				},
			},
		}})
	return flow
}

func (s *MessageFlow) DefaultMessage() models.Message {
	if s.defaultMsg == (models.Message{}) {
		return models.Message{Text: "what can i do for you?"}
	}

	return s.defaultMsg
}

func (s *MessageFlow) SetDefaultMessage(msg models.Message) *MessageFlow {
	s.defaultMsg = msg
	return s
}

func (s *MessageFlow) Add(expectMsg string, message models.Message) *MessageFlow {
	s.lockFlow.Lock()
	defer s.lockFlow.Unlock()
	s.flow[expectMsg] = message
	return s
}

func (s *MessageFlow) Get(expectMsg string) (models.Message, bool) {
	s.lockFlow.RLock()
	defer s.lockFlow.RUnlock()
	msg, ok := s.flow[expectMsg]
	return msg, ok
}

func (s *MessageFlow) IsEmpty() bool {
	s.lockFlow.RLock()
	defer s.lockFlow.RUnlock()
	return len(s.flow) == 0
}

func (s *MessageFlow) Buid(recipientID, inputMsg string) (models.SendMessageRequest, error) {

	if s.IsEmpty() {
		return models.SendMessageRequest{}, errors.New("empty flow")
	}

	message, exists := s.Get(inputMsg)

	if !exists {
		message = s.DefaultMessage()
	}

	return models.SendMessageRequest{
		MessagingType: "RESPONSE",
		RecipientID:   models.MessageRecipient{ID: recipientID},
		Message:       message,
	}, nil
}

type SimpleMessageUC struct {
	messageFlow  *MessageFlow
	postbackFlow *MessageFlow
}

func NewSimpleMessageUC() *SimpleMessageUC {
	return &SimpleMessageUC{messageFlow: NewMessageFlow(), postbackFlow: NewPostbackFlow()}
}

func (s *SimpleMessageUC) HandleWebHookRequest(r models.WehbookReq) error {
	if r.Object != "page" {
		return errors.New("unknown web hook object")
	}

	for _, we := range r.Entry {
		err := s.handleWebHookRequestEntry(we)
		if err != nil {
			return fmt.Errorf("handle webhook request entry: %w", err)
		}
	}

	return nil
}

func (s *SimpleMessageUC) handleWebHookRequestEntry(we models.Entry) error {

	var err error

	defer func() error {
		if err != nil {
			log.Println("handle message: %w", err)
			return fmt.Errorf("handle message: %w", err)
		}
		return nil
	}()

	if len(we.Messaging) == 0 {
		log.Println("there is no message entry")
		return errors.New("there is no message entry")
	}

	em := we.Messaging[0]

	if em.Postback != nil {
		return s.handlerPostback(em.Sender.ID, em.Postback)
	}

	if em.Message != nil {
		return s.handleMessage(em.Sender.ID, em.Message.Text)
	}

	return nil
}

func (s *SimpleMessageUC) handlerPostback(recipientID string, postbackReq *models.Postback) error {
	fmt.Println("POSTBACK")

	if s.postbackFlow == nil {
		return errors.New("postbackFlow cannot be nil")
	}

	msgText := postbackReq.Payload
	msgRequest, err := s.postbackFlow.Buid(recipientID, msgText)

	if err != nil {
		return fmt.Errorf("error building flow: %w", err)
	}

	return SendRespose(msgRequest)
}

func (s *SimpleMessageUC) handleMessage(recipientID, msgText string) error {

	if s.messageFlow == nil {
		return errors.New("messageFlow cannot be nil")
	}

	msgText = strings.TrimSpace(msgText)
	msgRequest, err := s.messageFlow.Buid(recipientID, msgText)

	if err != nil {
		return fmt.Errorf("error building flow: %w", err)
	}

	return SendRespose(msgRequest)

}

func SendRespose(msgRequest models.SendMessageRequest) error {

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("access_token", pageAccessToken).
		SetBody(msgRequest).
		Post("https://graph.facebook.com/v12.0/me/messages")

	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	log.Println("Message sent:", resp)
	return nil

}

func UploadImage(url string) (string, error) {
	client := resty.New()
	r, err := client.R().
		SetHeader("content-type", "application/json").
		SetHeader("Authorization", "Bearer "+pageAccessToken).
		SetBody(models.SendMessageRequest{
			Message: models.Message{
				Attachment: &models.Attachment{
					Type: "image",
					Payload: map[string]any{
						"url":         "https://t.ctcdn.com.br/FuwTWKvmTX8FiPEsg-Ou7MNEokc=/768x432/smart/i571145.jpeg",
						"is_reusable": true,
					},
				},
			},
		}).
		Post("https://graph.facebook.com/v12.0/me/message_attachments")

	if err != nil {
		return "", err
	}

	fmt.Println("resp", r, err)

	var response models.MediaAttachmentResponse

	if err = json.Unmarshal(r.Body(), &response); err != nil {
		return "", err
	}

	return response.AttachmentID, nil
}
