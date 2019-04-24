package model

import (
	"encoding/json"
	"fmt"
	"github.com/just1689/pyrate-server/db"
)

type Player struct {
	X        int
	Z        int
	Incoming chan []byte
	Outgoing chan []byte
	Stop     chan bool
	Send     chan []byte
}

func CreatePlayerAndStart(incoming, outgoing chan []byte) *Player {
	player := Player{}
	player.X = 50
	player.Z = 50
	player.Incoming = incoming
	player.Outgoing = outgoing
	player.Stop = make(chan bool) //??? to handle
	player.start()
	return &player
}

func (player Player) start() {
	go func() {
		select {
		case <-player.Stop:
			fmt.Println("Player loop stopping")
			return
		case b := <-player.Incoming:
			player.handleMessage(b)

		}

		//TODO:
	}()
}

func (player Player) handleMessage(b []byte) {
	m, err := bytesToMessage(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	//THIS NEEDS TO MOVE OUTSIDE OF package chat

	if m.Topic == "map-request" {

		body := MapRequestBody{}
		err := json.Unmarshal(m.Body, &body)
		if err != nil {
			fmt.Println(err)
			return
		}

		conn, err := db.Connect()
		if err != nil {
			fmt.Println(err)
			return
		}
		c := model.GetTilesChunkAsync(conn, body.X-100, body.X+100, body.Y-100, body.Y+100)
		count := 0
		for tile := range c {
			if tile.TileType == model.TileTypeWater {
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
			client.send <- mb
			count++
		}
		fmt.Println("Sent", count, "tiles")

	} else if m.Topic == "keyboard" {
		keyboard := KeyboardBody{}
		err := json.Unmarshal(m.Body, &keyboard)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
