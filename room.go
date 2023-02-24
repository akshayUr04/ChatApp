package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	//clients hold  all curent clients in the room
	clients map[*client]bool

	// join is a channel for clients that wishing to join the room
	join chan *client

	// leave is a channel for clients that wishing to leave the room
	leave chan *client

	//forward is a channel that holds incoming messeges that should be forwarded to other childs
	forward chan []byte
}

// newRoom create a new chat room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

const (
	socektBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socektBufferSize, WriteBufferSize: socektBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("servHTTP", err)
		return
	}
	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
