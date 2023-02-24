package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	//socekt is the web socket for the client
	socket *websocket.Conn

	//receive is a channel to recive message from other clients
	receive chan []byte

	//room is the room this client chating int
	room *room
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
