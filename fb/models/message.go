package models

type (
	WehbookReq struct {
		Object string  `json:"object"`
		Entry  []Entry `json:"entry"`
	}

	Entry struct {
		ID        string      `json:"id"`
		Time      int         `json:"time"`
		Messaging []Messaging `json:"messaging"`
	}

	Messaging struct {
		Timestamp int       `json:"timestamp"`
		Sender    Sender    `json:"sender"`
		Recipient Recipient `json:"recipient"`
		Message   *Message  `json:"message"`
		Postback  *Postback `json:"postback"`
	}

	Sender struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Recipient struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Postback struct {
		Mid       string `json:"mid"`
		Title     string `json:"title"`
		Payload   string `json:"payload"`
		Timestamp int    `json:"timestamp"`
	}

	Message struct {
		Mid        string      `json:"mid,omitempty"`
		Seq        int         `json:"seq,omitempty"`
		Text       string      `json:"text,omitempty"`
		Attachment *Attachment `json:"attachment,omitempty"`
	}

	Attachment struct {
		Type    string `json:"type"`
		Payload any    `json:"payload"`
	}

	PayloadButtons struct {
		TemplateType string   `json:"template_type"`
		Text         string   `json:"text"`
		Buttons      []Button `json:"buttons"`
	}

	Buttons []Button

	Button struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Payload string `json:"payload"`
	}

	PayloadCoupon struct {
		TemplateType         string `json:"template_type"`
		Title                string `json:"title"`
		Subtitle             string `json:"subtitle,omitempty"`
		CouponCode           string `json:"coupon_code"`
		CouponUrl            string `json:"coupon_url,omitempty"`
		CouponUrlButtonTitle string `json:"coupon_url_button_title,omitempty"`
		ImageUrl             string `json:"image_url,omitempty"`
		Payload              string `json:"payload"`
	}

	PayloadMedia struct {
		TemplateType string         `json:"template_type"`
		Elements     []MediaElement `json:"elements"`
	}

	MediaElement struct {
		MediaType    string  `json:"media_type"`
		AttachmentId string  `json:"attachment_id,omitempty"`
		Url          string  `json:"url"`
		Buttons      Buttons `json:"buttons,omitempty"`
	}

	MessageRecipient struct {
		ID string `json:"id"`
	}

	MediaAttachmentResponse struct {
		AttachmentID string `json:"attachment_id"`
	}

	SendMessageRequest struct {
		MessagingType string           `json:"messaging_type"`
		Tag           string           `json:"tag,omitempty"`
		RecipientID   MessageRecipient `json:"recipient"`
		Message       Message          `json:"message"`
	}
)
