package model

import (
	"encoding/json"
	"math"
)

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
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

func (o1 MessageOffset) equals(o2 MessageOffset) bool {
	return floatEquals(o1.X, o2.X) && floatEquals(o1.Y, o2.Y)
}

var eps float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if math.Abs(a-b) < eps {
		return true
	}
	//fmt.Println(math.Abs(a-b))
	return false
}
