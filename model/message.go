package model

// Message 消息体
type Message struct {
	Type   string `json:"type"`
	Data   string `json:"data"`
	Length int    `json:"length"`
	Sender string `json:"sender"`
}

// NewMessage 实例化
func NewMessage(t string, data string, length int, sender string) *Message {

	return &Message{
		Type:   t,
		Data:   data,
		Length: length,
		Sender: sender,
	}
}
