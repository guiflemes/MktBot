package service

import (
	"errors"
	"fmt"
	"log"
	"marketingBot/fb/adapters"
	"marketingBot/fb/models"
	"strings"

	dash "marketingBot/dashboard/adapters"
)

type SimpleMessageUC struct {
	messageFlow    *MessageFlow
	postbackFlow   *MessageFlow
	postbackAction []func(recipientID string, postback *models.Postback)
}

func NewSimpleMessageUC() *SimpleMessageUC {
	return &SimpleMessageUC{
		messageFlow:  NewMessageFlow(),
		postbackFlow: NewPostbackFlow(),
		postbackAction: []func(recipientID string, postback *models.Postback){
			collectFbButtonMetrics,
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

func (s *SimpleMessageUC) executePosbackAction(recipientID string, postback *models.Postback) {
	for _, fn := range s.postbackAction {
		go fn(recipientID, postback)
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

	s.executePosbackAction(recipientID, postbackReq)

	return adapters.SendRespose(msgRequest)
}

// TODO get coupon reveal on handleMessage

func (s *SimpleMessageUC) handleMessage(recipientID, msgText string) error {

	if s.messageFlow == nil {
		return errors.New("messageFlow cannot be nil")
	}

	msgText = strings.TrimSpace(msgText)
	msgRequest, err := s.messageFlow.Buid(recipientID, msgText)

	if err != nil {
		return fmt.Errorf("error building flow: %w", err)
	}

	return adapters.SendRespose(msgRequest)

}

func collectFbButtonMetrics(recipientID string, postbackReq *models.Postback) {
	repo := dash.ButtonStatisticsRepoMemory
	repo.Save(dash.ButtonClick{
		Title:      postbackReq.Title,
		Key:        postbackReq.Payload,
		Timestamp:  postbackReq.Timestamp,
		CustomerID: recipientID,
		Platform:   "FB",
	})
}
