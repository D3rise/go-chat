package main

import (
	"flag"
	"log"
	"net/http"

	"./structs"

	"golang.org/x/net/websocket"
)

func main() {
	flag.Parse()
	log.Println("Server listening on address: localhost:", *port)
	log.Fatal(server(*port))
}

var (
	port = flag.String("port", "9000", "port used for ws connection")
)

func server(port string) error {
	h := structs.NewHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handler(ws, h)
	}))

	s := http.Server{Addr: ":" + port, Handler: mux}
	return s.ListenAndServe()
}

func handler(ws *websocket.Conn, h *structs.Hub) {
	go h.Run()

	var u structs.User
	websocket.JSON.Receive(ws, &u)

	h.AddClientCh <- ws

	for {
		var m structs.Message
		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			h.RemoveClient(ws)
			return
		}
		h.BroadcastCh <- m
	}
}
