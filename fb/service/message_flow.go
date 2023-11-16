package service

import (
	"encoding/json"
	"fmt"
	"log"
	"marketingBot/dashboard/flow"
	"marketingBot/fb/models"
	"regexp"
	"sync"
)

func replacePlaceholders(templete string, replacements map[string]string) string {
	re := regexp.MustCompile(`{([^}]+)}`)
	return re.ReplaceAllStringFunc(templete, func(math string) string {
		key := math[1 : len(math)-1]
		if value, ok := replacements[key]; ok {
			return value
		}

		return math
	})
}

type MessageFlow struct {
	lockFlow   sync.RWMutex
	flow       map[string]func(sender models.Sender) models.Message
	defaultMsg func(sender models.Sender) models.Message
}

func (s *MessageFlow) DefaultMessageMaker() func(s models.Sender) models.Message {
	if s.defaultMsg == nil {
		return func(s models.Sender) models.Message {
			return models.Message{Text: "what can i do for you?"}
		}
	}

	return s.defaultMsg
}

func (s *MessageFlow) SetDefaultMessageMaker(maker func(sender models.Sender) models.Message) *MessageFlow {
	s.defaultMsg = maker
	return s
}

func (s *MessageFlow) Add(expectMsg string, messageMaker func(sender models.Sender) models.Message) *MessageFlow {
	s.lockFlow.Lock()
	defer s.lockFlow.Unlock()
	s.flow[expectMsg] = messageMaker
	return s
}

func (s *MessageFlow) Get(expectMsg string) (func(sender models.Sender) models.Message, bool) {
	s.lockFlow.RLock()
	defer s.lockFlow.RUnlock()

	maker, ok := s.flow[expectMsg]

	if !ok {
		return nil, false
	}
	return maker, true
}

func (s *MessageFlow) IsEmpty() bool {
	s.lockFlow.RLock()
	defer s.lockFlow.RUnlock()
	return len(s.flow) == 0
}

func (s *MessageFlow) Buid(sender models.Sender, inputMsg string) (models.SendMessageRequest, error) {

	messageMaker, exists := s.Get(inputMsg)

	if !exists {
		messageMaker = s.DefaultMessageMaker()
	}

	message := messageMaker(sender)

	return models.SendMessageRequest{
		MessagingType: "RESPONSE",
		RecipientID:   models.MessageRecipient{ID: sender.ID},
		Message:       message,
	}, nil
}

type BotFlow struct {
	postbackFlow *MessageFlow
	directFlow   *MessageFlow
}

func MessageFlowBuilder(flow *flow.Flow) (*BotFlow, error) {
	msgMaker := NewMessageMaker()
	msgMaker.SetMaker("button", ButtonMaker)
	msgMaker.SetMaker("image", ImageMaker)
	msgMaker.SetMaker("coupon", CouponMaker)

	flow.Lock()
	defer flow.Unlock()

	postbackFlow := &MessageFlow{flow: make(map[string]func(sender models.Sender) models.Message)}
	messageFlow := &MessageFlow{flow: make(map[string]func(sender models.Sender) models.Message)}

	for _, rel := range flow.Relationships {
		fmt.Println("postback")
		postback := flow.Cards[rel.TargetCardID]
		postbackMaker, err := msgMaker.Make(postback)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		postbackFlow.Add(rel.TargetCardID, postbackMaker)
		delete(flow.Cards, rel.TargetCardID)
	}

	for _, card := range flow.Cards {
		msgMaker, err := msgMaker.Make(card)

		fmt.Println("direcet")

		if err != nil {
			log.Println(err)
			return nil, err
		}

		if card.Initial {
			messageFlow.SetDefaultMessageMaker(msgMaker)
			continue
		}

		messageFlow.Add(card.ExpectedMsg, msgMaker)
	}

	return &BotFlow{postbackFlow: postbackFlow, directFlow: messageFlow}, nil

}

type messageMaker struct {
	makers map[string]func(template []byte) (func(models.Sender) models.Message, error)
	lock   sync.RWMutex
}

func NewMessageMaker() *messageMaker {
	return &messageMaker{
		makers: make(map[string]func(template []byte) (func(models.Sender) models.Message, error)),
	}
}

func (m *messageMaker) SetMaker(makeType string, maker func(template []byte) (func(models.Sender) models.Message, error)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.makers[makeType] = maker
}

func (m *messageMaker) Make(card flow.Card) (func(models.Sender) models.Message, error) {

	fmt.Println("CARD", card)
	m.lock.RLock()
	defer m.lock.RLock()

	template_json, err := json.Marshal(card.Template)
	if err != nil {
		log.Println("Failed to marshal template : ", err)
		return nil, fmt.Errorf("failed to marshal template : %w", err)
	}

	maker, ok := m.makers[card.Type]
	if !ok {
		log.Println("no card type found")
		return nil, fmt.Errorf("no card type %s found", card.Type)
	}

	return maker(template_json)
}

func CouponMaker(template []byte) (func(models.Sender) models.Message, error) {
	var coupon flow.CouponTemplate

	if err := json.Unmarshal(template, &coupon); err != nil {
		return nil, err
	}

	return func(s models.Sender) models.Message {
		return models.Message{Attachment: &models.Attachment{
			Type: "template",
			Payload: models.PayloadCoupon{
				TemplateType: "coupon",
				Title:        coupon.Title,
				Subtitle:     coupon.Subtitle,
				CouponCode:   coupon.Code,
				Payload:      coupon.Key,
			},
		}}
	}, nil
}

func ImageMaker(template []byte) (func(models.Sender) models.Message, error) {
	var image flow.ImageTemplate

	if err := json.Unmarshal(template, &image); err != nil {
		return nil, err
	}

	return func(s models.Sender) models.Message {
		return models.Message{Attachment: &models.Attachment{
			Type: "template",
			Payload: models.PayloadMedia{
				TemplateType: "media",
				Elements: []models.MediaElement{
					{
						MediaType:    "image",
						AttachmentId: image.ImageID,
					},
				},
			},
		}}
	}, nil
}

func ButtonMaker(template []byte) (func(models.Sender) models.Message, error) {
	var button flow.ButtonTemplate

	if err := json.Unmarshal(template, &button); err != nil {
		return nil, err
	}

	return func(s models.Sender) models.Message {
		var buttons []models.Button

		for _, option := range button.Options {
			optionPayload := models.OptionButtonPayload{
				TargetMessageID: option.TargetCardID,
				QuestionKey:     button.Key,
				OptionKey:       option.Key,
			}

			payload, _ := json.Marshal(optionPayload)
			buttons = append(buttons, models.Button{
				Type:    "postback",
				Title:   option.Text,
				Payload: string(payload),
			})
		}

		return models.Message{
			Attachment: &models.Attachment{
				Type: "template",
				Payload: models.PayloadButtons{
					TemplateType: "button",
					Text:         replacePlaceholders(button.Text, map[string]string{"name": s.Name}),
					Buttons:      buttons,
				},
			},
		}
	}, nil
}

func SampleBotFlowMock(key string) *BotFlow {
	mock := `{
		"name": "SampleFlow",
		"key": "sample_flow_key",
		"cards": {
		  "buttonCard1": {
			"id": "buttonCard1",
			"type": "button",
			"initial": true,
			"expected_msg": "",
			"template": {
			  "key": "welcome_demo",
			  "text": "Welcome to the demo promotional flow {name}! Are you interested in our coupon",
			  "options": [
				{"text": "Yes! Show me coupon", "target_card_id": "couponCard", "key": "yes"},
				{"text": "No, thanks", "target_card_id": "imageCard", "key": "no"}
			  ]
			}
		  },
		  "imageCard": {
			"id": "imageCard",
			"type": "image",
			"initial": false,
			"expected_msg": "",
			"template": {
			  "image_url": "",
			  "image_id": "1746931029114090"
			}
		  },
		  "couponCard": {
			"id": "couponCard",
			"type": "coupon",
			"initial": false,
			"expected_msg": "",
			"template": {
			  "title": "here is our unqiue promotinal coupon",
			  "subtitle": "10% off limit 1 per customer",
			  "code": "10FF",
			  "key": "10FF"
			}
		  }
		},
		"relationships": [
		  {
			"source_card_id": "buttonCard1",
			"target_card_id": "imageCard",
			"relationship_type": "button_to_image",
			"additional_details": "Option 1 selected"
		  },
		  {
			"source_card_id": "buttonCard1",
			"target_card_id": "couponCard",
			"relationship_type": "button_to_coupon",
			"additional_details": "Option 2 selected"
		  }
		]
	  }`

	var flow flow.Flow
	err := json.Unmarshal([]byte(mock), &flow)
	if err != nil {
		fmt.Println("Error:", err)
	}

	botFlow, err := MessageFlowBuilder(&flow)

	if err != nil {
		fmt.Println("Error:", err)
	}

	return botFlow

}
