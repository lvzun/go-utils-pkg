package notifySdk

type CallVoiceResponse struct {
	BaseModel
	Data string `json:"data"`
}

const (
	SEND_CHANNEL_QYWX  = "qywx"
	SEND_CHANNEL_SMS   = "sms"
	SEND_CHANNEL_EMAIL = "email"
	SEND_CHANNEL_WX    = "wx"
	SEND_CHANNEL_QQ    = "qq"
	SEND_CHANNEL_CALL  = "call"

	RECEIVER_TYPE_GROUP  = "group"
	RECEIVER_TYPE_PERSON = "person"
)

type SendMessageRequestParams struct {
	TaskId      string            `json:"task_id"`
	SendChannel string            `json:"send_channel"` //qywx phone sms email wx qq call
	Receiver    []*Receiver       `json:"receiver"`
	TemplateId  string            `json:"template_id"`
	Text        string            `json:"text"`
	Params      map[string]string `json:"params"`
}

type Receiver struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	QywxId       string `json:"qywxId"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	WxId         string `json:"wxId"`
	WxGroupId    string `json:"wxGroupId"`
	QqGroupId    int64  `json:"qqGroupId"`
	QqFriendId   int64  `json:"qqFriendId"`
	ReceiverType string `json:"receiverType"` //group person ,default:person
	ReceiverName string `json:"receiverName"`
	ReceiverId   string `json:"receiverId"`
	SendChannel  string `json:"send_channel"` //qywx phone sms email wx qq call
}
