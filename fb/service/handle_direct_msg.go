package service

import (
	"errors"
	"fmt"
	"log"
	"marketingBot/fb/models"
	"strings"
)

type DirectHandler struct {
	messageFlow *MessageFlow
	graphApi    GraphApiClientResponse
}

func (h *DirectHandler) Handle(sender models.Sender, msg *models.Message) error {
	if h.messageFlow == nil {
		return errors.New("messageFlow cannot be nil")
	}

	msgText := strings.TrimSpace(msg.Text)
	msgRequest, err := h.messageFlow.Buid(sender, msgText)

	if err != nil {
		log.Println("error building flow ", err)
		return fmt.Errorf("error building flow: %w", err)
	}

	return h.graphApi.SendRespose(msgRequest)

}
