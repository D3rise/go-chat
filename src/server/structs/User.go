package structs

import (
	"fmt"

	"golang.org/x/net/websocket"
)

// User struct
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// RemoveClient - Removes client from server connected clients list
func (h *Hub) RemoveClient(conn *websocket.Conn) {
	delete(h.Clients, conn.RemoteAddr().String())
	fmt.Println("Client disconnected.")
}

// AddClient - Add's connected client to server list of connected clients
func (h *Hub) AddClient(conn *websocket.Conn) {
	h.Clients[conn.RemoteAddr().String()] = conn
	fmt.Println("New client!")
}
