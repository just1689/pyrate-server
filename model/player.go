package model

import (
	"encoding/json"
	"fmt"
	"github.com/just1689/pyrate-server/db"
	"github.com/just1689/pyrate-server/queues"
	"math"
	"sync"
	"time"
)

const GlobalMaxPhysyicalDiff = 100

type Player struct {
	MidX      int
	MidZ      int
	Incoming  chan []byte
	Outgoing  chan []byte
	Stop      chan bool
	Send      chan []byte
	Keyboard  *KeyboardBody
	Offset    MessageOffset
	wgPhysics sync.WaitGroup

	lastOffset    MessageOffset
	NATSPublisher func(subject string, msg []byte) error
}

func CreatePlayerAndStart(incoming, outgoing chan []byte) *Player {
	player := Player{
		Incoming: incoming,
		Outgoing: outgoing,
		Stop:     make(chan bool), //??? to handle
		Offset: MessageOffset{
			X: -500,
			Z: -500,
		},
		lastOffset: MessageOffset{
			X: 0,
			Z: 0,
		},
		Keyboard: &KeyboardBody{},
	}
	player.MidX = math.Ilogb(math.Abs(player.Offset.X))
	player.MidZ = math.Ilogb(math.Abs(player.Offset.Z)) //TODO: no nonsense

	var err error
	player.NATSPublisher, err = queues.GetNATSPublisher()
	if err != nil {
		fmt.Println(err)
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

	go func() {
		for {
			time.Sleep(1 * time.Second)
			player.checkForMapDiff()
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
	body.MakeAbs()

	conn, err := db.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	c := GetTilesChunkAsync(conn, body.X-100, body.X+100, body.Z-100, body.Z+100)
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
		player.Offset.Z -= 0.1
	}
	if player.Keyboard.S {
		player.Offset.Z += 0.1
	}
}

func (player *Player) sendOffset() {
	if !player.Offset.equals(player.lastOffset) {

		player.lastOffset.X = player.Offset.X
		player.lastOffset.Z = player.Offset.Z

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

func (player *Player) checkForMapDiff() bool {

	tX := math.Ilogb(math.Abs(player.Offset.X))
	if tX-player.MidX > GlobalMaxPhysyicalDiff {
		player.MidX = tX
		return true
	}

	tY := math.Ilogb(math.Abs(player.Offset.Z))
	if tY-player.MidZ > GlobalMaxPhysyicalDiff {
		player.MidX = tX
		return true
	}

	return false

}
