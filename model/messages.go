package model

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
	W bool `json:"W"`
	A bool `json:"A"`
	S bool `json:"S"`
	D bool `json:"D"`
}

type MessageOffset struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
}
