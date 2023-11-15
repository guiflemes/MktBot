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
		Timestamp int                `json:"timestamp"`
		Sender    Sender             `json:"sender"`
		Recipient Recipient          `json:"recipient"`
		Message   *Message           `json:"message"`
		Postback  *Postback          `json:"postback"`
		Template  *MessagingTemplate `json:"template,omitempty"`
	}

	MessagingTemplate struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
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

	OptionButtonPayload struct {
		TargetMessageID string `json:"target_message_id"`
		QuestionKey     string `json:"question_key"`
		OptionKey       string `json:"option_key"`
	}
)

// {"object":"page","entry":[{"id":"171792666015676","time":1700062750516,"messaging":[{"sender":{"id":"171792666015676"},"recipient":{"id":"24536016889377737"},"timestamp":1700062750062,"message":{"mid":"m_8j_CHfYQ3voeCvGsKVf3RleFKihbgyHO3zZR-ByjmL16ypwiOleqgcedB-Xn-CU0U7xB3DbSTHM5GFjPFToM2A","is_echo":true,"app_id":866644431628954,"attachments":[{"type":"template","title":"here is our unqiue promotinal coupon","url":null,"payload":{"template_type":"generic","sharable":false,"elements":[{"title":"here is our unqiue promotinal coupon","subtitle":"10\u0025 off limit 1 per customer","image_url":null,"buttons":[{"page_id":171792666015676,"app_id":866644431628954,"payload":"10off","coupon_code":"10FF"}],"default_action":null}]}}]}}]}]}

// reveal cooupom
//{"object":"page","entry":[{"id":"171792666015676","time":1700062778479,"messaging":[{"sender":{"id":"24536016889377737"},"recipient":{"id":"171792666015676"},"timestamp":1700062777827,"template":{"type":"coupon","coupon_code":"10FF","payload":"10off"}}]}]}
