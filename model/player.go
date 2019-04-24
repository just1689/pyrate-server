package model

import "fmt"

type Player struct {
	X        int
	Z        int
	Messages chan []byte
	Stop     chan bool
}

func CreatePlayer() *Player {
	player := Player{}
	player.X = 50
	player.Z = 50
	player.Messages = make(chan []byte)
	player.Stop = make(chan bool)
}

func (player Player) Start() {
	go func() {
		select {
		case <-player.Stop:
			fmt.Println("Player loop stopping")
			return
		case b := <-player.Messages:
			player.handleMessage(b)

		}

		//TODO:
	}()
}

func (player Player) handleMessage(bytes []byte) {
	//
}
