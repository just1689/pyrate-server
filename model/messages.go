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
	Z int `json:"Z"`
}

func (m *MapRequestBody) MakeAbs() {
	if m.X < 0 {
		m.X /= -1
	}
	if m.Z < 0 {
		m.Z /= -1
	}
}

type KeyboardBody struct {
	W bool `json:"W"`
	A bool `json:"A"`
	S bool `json:"S"`
	D bool `json:"D"`
}

type MessageOffset struct {
	X float64 `json:"X"`
	Z float64 `json:"Z"`
}

func (o1 MessageOffset) equals(o2 MessageOffset) bool {
	return floatEquals(o1.X, o2.X) && floatEquals(o1.Z, o2.Z)
}

var eps float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if math.Abs(a-b) < eps {
		return true
	}
	//fmt.Println(math.Abs(a-b))
	return false
}
