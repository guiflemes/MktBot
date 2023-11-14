package adapters

import (
	"encoding/json"
	"fmt"
	"log"
	"marketingBot/fb/models"
	"marketingBot/settings"

	"github.com/go-resty/resty/v2"
)

func SendRespose(msgRequest models.SendMessageRequest) error {

	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("access_token", settings.GETENV("FB_PAGE_ACCESS_TOKEN")).
		SetBody(msgRequest).
		Post("https://graph.facebook.com/v12.0/me/messages")

	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	return nil

}

func UploadImage(url string) (string, error) {
	client := resty.New()
	r, err := client.R().
		SetHeader("content-type", "application/json").
		SetHeader("Authorization", "Bearer "+settings.GETENV("FB_PAGE_ACCESS_TOKEN")).
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
