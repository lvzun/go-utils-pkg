package notify

type Receiver struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	ReceiverName string `json:"receiverName"`
	ReceiverId   string `json:"receiverId"`
	ReceiverType string `json:"receiverType"`
	SendChannel  string `json:"send_channel"`
}

type SendMessageApiRequestDto struct {
	TaskId     string            `json:"task_id"`
	Receiver   []*Receiver       `json:"receiver"`
	TemplateId string            `json:"template_id"`
	Text       string            `json:"text"`
	Params     map[string]string `json:"params"`
}
type SendMessageApiResponseDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
