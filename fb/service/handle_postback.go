package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"marketingBot/fb/models"
)

type PostbackHandler struct {
	postbackAction []func(recipientID string, postback PostBackMetric)
	postbackFlow   *MessageFlow
	graphApi       GraphApiClientResponse
}

func (h *PostbackHandler) executePosbackAction(sender models.Sender, postback PostBackMetric) {
	for _, fn := range h.postbackAction {
		go fn(sender.ID, postback)
	}
}

func (h *PostbackHandler) Handle(sender models.Sender, postbackReq *models.Postback) error {
	if h.postbackFlow == nil {
		return errors.New("postbackFlow cannot be nil")
	}

	var option models.OptionButtonPayload

	if err := json.Unmarshal([]byte(postbackReq.Payload), &option); err != nil {
		log.Println("failed to unmarshal postback Payload : ", err)
		return fmt.Errorf("failed to unmarshal postback Payload : %w", err)
	}

	msgRequest, err := h.postbackFlow.Buid(sender, option.TargetMessageID)

	if err != nil {
		return fmt.Errorf("error building flow: %w", err)
	}

	h.executePosbackAction(sender, PostBackMetric{
		Title:     postbackReq.Title,
		Payload:   option,
		Timestamp: postbackReq.Timestamp,
	})

	return h.graphApi.SendRespose(msgRequest)
}
