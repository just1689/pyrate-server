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

type MapRequestBody struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type KeyboardBody struct {
	W int `json:"W"`
	A int `json:"A"`
	S int `json:"S"`
	D int `json:"D"`
}
