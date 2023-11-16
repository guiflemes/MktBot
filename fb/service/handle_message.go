package service

import (
	"errors"
	"fmt"
	"log"
	"marketingBot/fb/adapters"
	"marketingBot/fb/models"
	"time"

	dash "marketingBot/dashboard/adapters"
)

type (
	HandlerDirectMsg interface {
		Handle(sender models.Sender, msg *models.Message) error
		SetFlow(messageFlow *MessageFlow)
	}

	HandlerPostBackMsg interface {
		Handle(sender models.Sender, postbackReq *models.Postback) error
		SetFlow(messageFlow *MessageFlow)
	}

	HandleTemplateMsg interface {
		Handler(sender models.Sender, temp *models.MessagingTemplate)
	}

	SenderCacher interface {
		GetSenderName(senderID string, fetchGeSanderName func(sendId string) (string, error)) (string, error)
	}

	GraphApiSender interface {
		GetSenderName(senderID string) (string, error)
	}

	GraphApiClientResponse interface {
		SendRespose(msgRequest models.SendMessageRequest) error
	}
)

func MemoryGetBotFlow(key string) *BotFlow {
	repo := dash.MemoryFlowRepo
	flow := repo.Get(key)

	if flow == nil {
		return nil
	}

	fmt.Println(flow.Relationships)
	f, err := MessageFlowBuilder(flow)
	if err != nil {
		log.Println("error building a message flow")
		return nil
	}

	return f
}

type SimpleMessageUC struct {
	senderCache     SenderCacher
	graphApi        GraphApiSender
	templateHandler HandleTemplateMsg
	postbackHandler HandlerPostBackMsg
	directHandler   HandlerDirectMsg
	getBotFlow      func(key string) *BotFlow
}

func NewSimpleMessageUC() *SimpleMessageUC {

	cache := adapters.NewSenderCache()
	graphAPi := adapters.NewGrapApi()

	return &SimpleMessageUC{
		senderCache: cache,
		graphApi:    graphAPi,
		templateHandler: &TemplateHandler{templateAction: []func(recipientID string, template *models.MessagingTemplate){
			collecFbCoupomRevelMetric,
		}},
		postbackHandler: &PostbackHandler{
			postbackAction: []func(recipientID string, postback PostBackMetric){
				collectFbButtonMetrics,
			},
			graphApi: graphAPi,
		},
		directHandler: &DirectHandler{
			graphApi: graphAPi,
		},
		getBotFlow: SampleBotFlowMock,
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

func (s *SimpleMessageUC) handleWebHookRequestEntry(we models.Entry) error {

	flow := s.getBotFlow("sample_flow_key")

	if flow == nil {
		return fmt.Errorf("bot flow with thje given key '%s' not found", "sample_bot_flow")
	}

	s.postbackHandler.SetFlow(flow.postbackFlow)
	s.directHandler.SetFlow(flow.directFlow)

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
		return s.postbackHandler.Handle(m.Sender, m.Postback)
	}

	if m.Message != nil {
		return s.directHandler.Handle(m.Sender, m.Message)
	}

	if m.Template != nil {
		s.templateHandler.Handler(m.Sender, m.Template)
	}

	return nil
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
