package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"./structs"

	"golang.org/x/net/websocket"
)

func main() {
	flag.Parse()

	ws, err := connect("https://localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected")

	var m structs.Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error receiving message: ", err.Error())
				break
			}

			log.Printf("<%s>: \"%s\"\n", m.Author, m.Text)
			fmt.Print("> ")
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	login := scanner.Text()
	fmt.Print("\nPassword: ")
	scanner.Scan()
	password := scanner.Text()

	user := structs.User{
		Login:    login,
		Password: password,
	}

	websocket.JSON.Send(ws, user)
	fmt.Print("> ")

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		m := structs.Message{
			Text:   text,
			Author: login,
		}

		err = websocket.JSON.Send(ws, m)

		if err != nil {
			fmt.Println("Error sending message: ", err.Error())
			break
		}
	}
}

var (
	port = flag.String("port", "9000", "port used for ws connection")
)

func connect(address string) (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", mockedIP())
}

func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}

	return fmt.Sprintf("http://%d:%d:%d:%d", arr[0], arr[1], arr[2], arr[3])
}
