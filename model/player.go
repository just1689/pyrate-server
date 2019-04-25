package model

import (
	"encoding/json"
	"fmt"
	"github.com/just1689/pyrate-server/db"
	"sync"
	"time"
)

type Player struct {
	X         int
	Z         int
	Incoming  chan []byte
	Outgoing  chan []byte
	Stop      chan bool
	Send      chan []byte
	Keyboard  *KeyboardBody
	Offset    MessageOffset
	wgPhysics sync.WaitGroup

	lastOffset MessageOffset
}

func CreatePlayerAndStart(incoming, outgoing chan []byte) *Player {
	player := Player{
		X:        50,
		Z:        50,
		Incoming: incoming,
		Outgoing: outgoing,
		Stop:     make(chan bool), //??? to handle
		Offset: MessageOffset{
			X: -25,
			Y: -25,
		},
		lastOffset: MessageOffset{
			X: 0,
			Y: 0,
		},
		Keyboard: &KeyboardBody{},
	}

	player.wgPhysics.Add(1)
	player.start()
	return &player
}

func (player *Player) start() {

	//Message handler
	go func() {
		for {
			select {
			case <-player.Stop:
				fmt.Println("Player loop stopping")
				return
			case b := <-player.Incoming:
				player.handleMessage(b)
			case <-time.After(60 * time.Millisecond):
				player.sendOffset()
			}

		}
	}()

	//Physics handler
	go func() {
		for {
			select {
			case <-time.After(60 * time.Millisecond):
				player.move()

			}
		}
	}()

}

func (player *Player) handleMessage(b []byte) {
	m, err := bytesToMessage(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	//THIS NEEDS TO MOVE OUTSIDE OF package chat

	if m.Topic == "map-request" {
		player.handleMapRequest(m.Body)

	} else if m.Topic == "keyboard" {
		player.handleKeyboardRequest(m.Body)

	}
}

func (player *Player) handleMapRequest(rawBody json.RawMessage) {
	body := MapRequestBody{}
	err := json.Unmarshal(rawBody, &body)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := db.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	c := GetTilesChunkAsync(conn, body.X-100, body.X+100, body.Y-100, body.Y+100)
	count := 0
	for tile := range c {
		if tile.TileType == TileTypeWater {
			continue
		}
		b, err := json.Marshal(*tile)
		m := Message{
			Topic: "tile",
			Body:  b,
		}
		mb, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return
		}
		player.Outgoing <- mb
		count++
	}
	fmt.Println("Sent", count, "tiles")

}

func (player *Player) handleKeyboardRequest(messages json.RawMessage) {
	keyboard := KeyboardBody{}
	err := json.Unmarshal(messages, &keyboard)
	if err != nil {
		fmt.Println(err)
		return
	}
	player.Keyboard = &keyboard

}

func (player *Player) move() {
	if player.Keyboard.A {
		player.Offset.X += 0.1
	}
	if player.Keyboard.D {
		player.Offset.X -= 0.1
	}
	if player.Keyboard.W {
		player.Offset.Y -= 0.1
	}
	if player.Keyboard.S {
		player.Offset.Y += 0.1
	}
}

func (player *Player) sendOffset() {
	if !player.Offset.equals(player.lastOffset) {

		player.lastOffset.X = player.Offset.X
		player.lastOffset.Y = player.Offset.Y

		bo, err := json.Marshal(player.Offset)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := Message{
			Topic: "offset",
			Body:  bo,
		}

		b, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		player.Outgoing <- b

	}

}
