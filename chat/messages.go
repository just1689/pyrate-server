package chat

import "encoding/json"

type Message struct {
	Topic string          `json:"topic"`
	Body  json.RawMessage `json:"body"`
}

func bytesToMessage(b []byte) (m *Message, err error) {
	m = &Message{}
	err = json.Unmarshal(b, m)
	return
}
