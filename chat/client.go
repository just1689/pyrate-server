package chat

import (
	"bytes"
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
	SendToWS      chan []byte
	StopNSQ       chan bool
	StopPlayer    chan bool
	SendToPlayer  chan []byte
}

func (c *Client) Auth() bool {
	ok := true
	c.Authenticated = ok
	return ok
}

func (c *Client) close() {
	c.StopNSQ <- true
	c.StopPlayer <- true

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
		c.SendToPlayer <- message
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
		case message, ok := <-c.SendToWS:
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

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request, name string, secret string, subscriber func(topic, channel string) chan bool, PlayerCreator func(c *Client)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:          hub,
		conn:         conn,
		SendToWS:     make(chan []byte, 256),
		id:           name,
		secret:       secret,
		StopPlayer:   make(chan bool),
		SendToPlayer: make(chan []byte, 256),
	}
	PlayerCreator(client)

	client.Auth()
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	if client.Authenticated {
		//client.StopNSQ = subscriber(fmt.Sprint("player."+client.id), "all")
	}

}
