package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	flag.Parse()
	log.Fatal(server(*port))
}

// Message struct
type Message struct {
	Text string `json:"text"`
}

type hub struct {
	clients        map[string]*websocket.Conn
	addClientCh    chan *websocket.Conn
	removeClientCh chan *websocket.Conn
	broadcastCh    chan Message
}

var (
	port = flag.String("port", "9000", "port used for ws connection")
)

func server(port string) error {
	h := newHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handler(ws, h)
	}))

	s := http.Server{Addr: ":" + port, Handler: mux}
	return s.ListenAndServe()
}

func handler(ws *websocket.Conn, h *hub) {
	go h.run()

	h.addClientCh <- ws

	for {
		var m Message
		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			h.broadcastCh <- Message{err.Error()}
			h.removeClient(ws)
			return
		}
		h.broadcastCh <- m
	}
}

func newHub() *hub {
	return &hub{
		clients:        make(map[string]*websocket.Conn),
		addClientCh:    make(chan *websocket.Conn),
		removeClientCh: make(chan *websocket.Conn),
		broadcastCh:    make(chan Message),
	}
}

func (h *hub) run() {
	for {
		select {
		case conn := <-h.addClientCh:
			h.addClient(conn)

		case conn := <-h.removeClientCh:
			h.removeClient(conn)

		case m := <-h.broadcastCh:
			h.broadcastMessage(m)
		}
	}
}

func (h *hub) broadcastMessage(m Message) {
	for _, conn := range h.clients {
		err := websocket.JSON.Send(conn, m)

		fmt.Printf("New message: %s", m)

		if err != nil {
			fmt.Printf("Error broadcasting message: ", err)
			return
		}
	}
}

func (h *hub) removeClient(conn *websocket.Conn) {
	delete(h.clients, conn.LocalAddr().String())
	fmt.Println("Client disconnected.")
}
func (h *hub) addClient(conn *websocket.Conn) {
	h.clients[conn.RemoteAddr().String()] = conn
	fmt.Println("New client!")
}
