package service

import (
	"marketingBot/fb/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageFlowDefaultMessage(t *testing.T) {

	assert := assert.New(t)
	flow1 := &MessageFlow{flow: make(map[string]func(sender models.Sender) models.Message)}
	flow2 := &MessageFlow{flow: make(map[string]func(sender models.Sender) models.Message)}

	maker := func(sender models.Sender) models.Message {
		return models.Message{Text: "test1"}
	}
	flow2.SetDefaultMessageMaker(maker)

	type testCase struct {
		desc     string
		expected func(sender models.Sender) models.Message
		Flow     *MessageFlow
	}

	for _, c := range []testCase{
		{
			desc:     "flow with default set pre defined",
			expected: flow1.DefaultMessageMaker(),
			Flow:     flow1,
		},
		{
			desc:     "flow with set default maker",
			expected: maker,
			Flow:     flow2,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			expectedMsg := c.expected(models.Sender{})
			fn := c.Flow.DefaultMessageMaker()
			msg := fn(models.Sender{})
			assert.Equal(expectedMsg.Text, msg.Text)
		})
	}

}

func TestMessageFlowBuilde(t *testing.T) {

	assert := assert.New(t)

	flow1 := &MessageFlow{flow: make(map[string]func(sender models.Sender) models.Message)}
	flow1.Add("test1", func(sender models.Sender) models.Message {
		return models.Message{Text: "expected test1"}
	})
	flow1.Add("test2", func(sender models.Sender) models.Message {
		return models.Message{Text: "expected test2"}
	})

	type testCase struct {
		desc      string
		msgInput  string
		expectMsg string
	}

	for _, c := range []testCase{
		{
			desc:      "retuns a message to test1",
			msgInput:  "test1",
			expectMsg: "expected test1",
		},
		{
			desc:      "retuns a message to test2",
			msgInput:  "test2",
			expectMsg: "expected test2",
		},
		{
			desc:      "return default msg",
			msgInput:  "test3",
			expectMsg: "what can i do for you?",
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			req, _ := flow1.Buid(models.Sender{}, c.msgInput)
			assert.Equal(req.Message.Text, c.expectMsg)
		})
	}

}

func TestMessageFlowGet(t *testing.T) {
	assert := assert.New(t)

	flow1 := &MessageFlow{flow: make(map[string]func(sender models.Sender) models.Message)}

	flow1.Add("test1", func(sender models.Sender) models.Message {
		return models.Message{Text: "expected test1"}
	})

	type testCase struct {
		desc       string
		msgInput   string
		expectBool bool
	}

	for _, c := range []testCase{
		{
			desc:       "registered key test1",
			msgInput:   "test1",
			expectBool: true,
		},
		{
			desc:       "not registered key test3",
			msgInput:   "test3",
			expectBool: false,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			_, ok := flow1.Get(c.msgInput)
			assert.Equal(ok, c.expectBool)
		})
	}

}
