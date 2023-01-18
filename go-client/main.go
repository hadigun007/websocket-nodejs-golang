package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	connectToNodejsWebsocketServer()
	fmt.Println("closing app")
}

// Client used for websocket
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 65536
)

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
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
		fmt.Println(string(message))
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	}
}

func connectToNodejsWebsocketServer() {
	url := "127.0.0.1:8080/foo"
	ws, res, err := websocket.DefaultDialer.Dial("ws://"+url, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	client := &Client{conn: ws, send: make(chan []byte, 256)}

	go client.readPump()

	for {
		// infinite loop so our client can wait for read/write from/to server
	}
}
