package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/just1689/pyrate-server/db"
	"github.com/just1689/pyrate-server/model"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	id            string
	secret        string
	Authenticated bool
	hub           *Hub
	conn          *websocket.Conn
	send          chan []byte
	StopNSQ       chan bool
}

func (c *Client) Auth() bool {
	ok := true
	c.Authenticated = ok
	return ok
}

func (c *Client) close() {
	c.StopNSQ <- true

}

func (c *Client) Check(in string) bool {
	return c.Auth()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()

		//Close the nsq thing
		c.close()

	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		messageHandler(c, message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request, name string, secret string, subscriber func(topic, channel string) chan bool) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		id:     name,
		secret: secret,
	}
	client.Auth()
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	if client.Authenticated {
		//client.StopNSQ = subscriber(fmt.Sprint("player."+client.id), "all")
	}

}

func messageHandler(client *Client, b []byte) {
	m, err := bytesToMessage(b)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	}

}
