package adapters

import (
	"encoding/json"
	"fmt"
	"log"
	"marketingBot/fb/models"
	"marketingBot/settings"

	"github.com/go-resty/resty/v2"
)

type GrapApi struct {
	client       *resty.Client
	baseUrl      string
	access_token string
}

func NewGrapApi() *GrapApi {
	return &GrapApi{
		client:       resty.New(),
		baseUrl:      "https://graph.facebook.com",
		access_token: settings.GETENV("FB_PAGE_ACCESS_TOKEN"),
	}
}

func (g *GrapApi) GetSenderName(senderID string) (string, error) {
	url := fmt.Sprintf("%s/v13.0/%s?fields=first_name,last_name&access_token=%s", g.baseUrl, senderID, g.access_token)

	resp, err := g.client.R().Get(url)

	if err != nil {
		return "", err
	}

	var profileData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &profileData); err != nil {
		return "", err
	}

	firstName, _ := profileData["first_name"].(string)

	return firstName, nil
}

func (g *GrapApi) UploadImage(url string) (string, error) {
	resp, err := g.client.R().
		SetHeader("content-type", "application/json").
		SetHeader("Authorization", "Bearer "+g.access_token).
		SetBody(models.SendMessageRequest{
			Message: models.Message{
				Attachment: &models.Attachment{
					Type: "image",
					Payload: map[string]any{
						"url":         url,
						"is_reusable": true,
					},
				},
			},
		}).
		Post(fmt.Sprintf("%s/v12.0/me/message_attachments", g.baseUrl))

	if err != nil {
		return "", err
	}

	var response models.MediaAttachmentResponse

	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return "", err
	}

	return response.AttachmentID, nil
}

func (g *GrapApi) SendRespose(msgRequest models.SendMessageRequest) error {
	_, err := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("access_token", g.access_token).
		SetBody(msgRequest).
		Post(fmt.Sprintf("%s/v12.0/me/messages", g.baseUrl))

	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}
	return nil

}
