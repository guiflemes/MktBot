package service

import (
	"marketingBot/fb/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type GraphApiStub struct {
	msgReq models.SendMessageRequest
}

func (g *GraphApiStub) SendRespose(msgRequest models.SendMessageRequest) error {
	g.msgReq = msgRequest
	return nil
}

func TestMessageFlowDefaultMessageNullMessageFlow(t *testing.T) {
	assert := assert.New(t)
	handler := &DirectHandler{}
	result := handler.Handle(models.Sender{}, &models.Message{})
	assert.EqualError(result, "messageFlow cannot be nil")
}

func TestMessageFlowDefaultMessageWithMessageFlow(t *testing.T) {
	assert := assert.New(t)
	flow := &MessageFlow{flow: make(map[string]func(sender models.Sender) models.Message)}
	flow.Add("test1", func(sender models.Sender) models.Message {
		return models.Message{
			Text: "testing",
		}
	})

	graphApi := &GraphApiStub{}

	handler := &DirectHandler{
		messageFlow: flow,
		graphApi:    graphApi,
	}
	result := handler.Handle(models.Sender{}, &models.Message{Text: "test1"})
	assert.NoError(result)
	assert.NotNil(graphApi.msgReq)

}
