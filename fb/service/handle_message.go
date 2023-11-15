package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"marketingBot/fb/adapters"
	"marketingBot/fb/models"
	"strings"
	"time"

	dash "marketingBot/dashboard/adapters"
)

type SenderCacher interface {
	GetSenderName(senderID string, fetchGeSanderName func(sendId string) (string, error)) (string, error)
}
type GraphApiClient interface {
	GetSenderName(senderID string) (string, error)
	SendRespose(msgRequest models.SendMessageRequest) error
}

type SimpleMessageUC struct {
	messageFlow    *MessageFlow
	postbackFlow   *MessageFlow
	postbackAction []func(recipientID string, postback PostBackMetric)
	senderCache    SenderCacher
	graphApi       GraphApiClient
	templateAction []func(recipientID string, template *models.MessagingTemplate)
}

func NewSimpleMessageUC() *SimpleMessageUC {

	flow := SampleBotFlowMock()
	cache := adapters.NewSenderCache()

	return &SimpleMessageUC{
		messageFlow:  flow.directFlow,
		postbackFlow: flow.postbackFlow,
		postbackAction: []func(recipientID string, postback PostBackMetric){
			collectFbButtonMetrics,
		},
		senderCache: cache,
		graphApi:    adapters.NewGrapApi(),
		templateAction: []func(recipientID string, temp *models.MessagingTemplate){
			collecFbCoupomRevelMetric,
		},
	}
}

func (s *SimpleMessageUC) HandleWebHookRequest(r models.WehbookReq) error {
	if r.Object != "page" {
		return errors.New("unknown webhook object")
	}

	for _, we := range r.Entry {
		err := s.handleWebHookRequestEntry(we)
		if err != nil {
			return fmt.Errorf("handle webhook request entry: %w", err)
		}
	}

	return nil
}

func (s *SimpleMessageUC) executePosbackAction(sender models.Sender, postback PostBackMetric) {
	for _, fn := range s.postbackAction {
		go fn(sender.ID, postback)
	}
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

	m := we.Messaging[0]
	senderName, err := s.senderCache.GetSenderName(m.Sender.ID, s.graphApi.GetSenderName)

	if err != nil {
		log.Println("failed getting send name: ", err)
		return fmt.Errorf("failed getting send name: %w", err)
	}

	m.Sender.Name = senderName

	if m.Postback != nil {
		return s.handlerPostback(m.Sender, m.Postback)
	}

	if m.Message != nil {
		return s.handleMessage(m.Sender, m.Message)
	}

	if m.Template != nil {
		s.handlerTemplate(m.Sender, m.Template)
	}

	return nil
}

func (s *SimpleMessageUC) handlerTemplate(sender models.Sender, temp *models.MessagingTemplate) {
	for _, fn := range s.templateAction {
		go fn(sender.ID, temp)
	}
}

func (s *SimpleMessageUC) handlerPostback(sender models.Sender, postbackReq *models.Postback) error {
	if s.postbackFlow == nil {
		return errors.New("postbackFlow cannot be nil")
	}

	var option models.OptionButtonPayload

	if err := json.Unmarshal([]byte(postbackReq.Payload), &option); err != nil {
		log.Println("failed to unmarshal postback Payload : ", err)
		return fmt.Errorf("failed to unmarshal postback Payload : %w", err)
	}

	msgRequest, err := s.postbackFlow.Buid(sender, option.TargetMessageID)

	if err != nil {
		return fmt.Errorf("error building flow: %w", err)
	}

	s.executePosbackAction(sender, PostBackMetric{
		Title:     postbackReq.Title,
		Payload:   option,
		Timestamp: postbackReq.Timestamp,
	})

	return s.graphApi.SendRespose(msgRequest)
}

func (s *SimpleMessageUC) handleMessage(sender models.Sender, msg *models.Message) error {
	if s.messageFlow == nil {
		return errors.New("messageFlow cannot be nil")
	}

	msgText := strings.TrimSpace(msg.Text)
	msgRequest, err := s.messageFlow.Buid(sender, msgText)

	if err != nil {
		log.Println("error building flow ", err)
		return fmt.Errorf("error building flow: %w", err)
	}

	return s.graphApi.SendRespose(msgRequest)

}

type TemplateHandler struct {
	templateAction []func(recipientID string, template *models.MessagingTemplate)
}

func (h *TemplateHandler) Handler(sender models.Sender, temp *models.MessagingTemplate) {
	for _, fn := range h.templateAction {
		go fn(sender.ID, temp)
	}
}

type MessageHandler struct {
	messageFlow *MessageFlow
	graphApi    GraphApiClient
}

func (h *MessageHandler) Handle(sender models.Sender, msg *models.Message) error {
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

type PostbackHandler struct {
	postbackAction []func(recipientID string, postback PostBackMetric)
	postbackFlow   *MessageFlow
	graphApi       GraphApiClient
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

type PostBackMetric struct {
	Title     string
	Timestamp int
	Payload   models.OptionButtonPayload
}

func collectFbButtonMetrics(recipientID string, postbackReq PostBackMetric) {
	repo := dash.StatisticsRepoMemory
	repo.SaveClicks(dash.ButtonClick{
		Title:       postbackReq.Title,
		QuestionKey: postbackReq.Payload.QuestionKey,
		OptionKey:   postbackReq.Payload.OptionKey,
		Timestamp:   postbackReq.Timestamp,
		CustomerID:  recipientID,
		Platform:    "FB",
	})
}

func collecFbCoupomRevelMetric(recipientID string, template *models.MessagingTemplate) {

	if template.Type != "coupon" {
		return
	}

	repo := dash.StatisticsRepoMemory
	repo.SaveRevels(dash.CouponRevel{Code: template.Payload, CustomerID: recipientID, Platform: "FB", Timestamp: time.Now().Unix()})
}
