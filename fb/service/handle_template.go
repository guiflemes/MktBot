package service

import "marketingBot/fb/models"

type TemplateHandler struct {
	templateAction []func(recipientID string, template *models.MessagingTemplate)
}

func (h *TemplateHandler) Handler(sender models.Sender, temp *models.MessagingTemplate) {
	for _, fn := range h.templateAction {
		go fn(sender.ID, temp)
	}
}
