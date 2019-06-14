package structs

import (
	"fmt"

	"golang.org/x/net/websocket"
)

// Hub - structure of the hub
type Hub struct {
	Clients        map[string]*websocket.Conn
	AddClientCh    chan *websocket.Conn
	RemoveClientCh chan *websocket.Conn
	BroadcastCh    chan Message
}

// NewHub - generate new hub
func NewHub() *Hub {
	return &Hub{
		Clients:        make(map[string]*websocket.Conn),
		AddClientCh:    make(chan *websocket.Conn),
		RemoveClientCh: make(chan *websocket.Conn),
		BroadcastCh:    make(chan Message),
	}
}

// Run - starts hub
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.AddClientCh:
			h.AddClient(conn)

		case conn := <-h.RemoveClientCh:
			h.RemoveClient(conn)

		case m := <-h.BroadcastCh:
			h.BroadcastMessage(m)
		}
	}
}

// BroadcastMessage - send message to all clients connected to server
func (h *Hub) BroadcastMessage(m Message) {
	for _, conn := range h.Clients {
		err := websocket.JSON.Send(conn, m)

		fmt.Printf("New message from \"%s\": %s\n", m.Author, m.Text)

		if err != nil {
			fmt.Printf("Error broadcasting message: %s\n", err)
			return
		}
	}
}
