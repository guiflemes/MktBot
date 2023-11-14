package service

import (
	"errors"
	"marketingBot/fb/models"
	"sync"
)

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
