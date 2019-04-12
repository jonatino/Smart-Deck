package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{}

var webhookUrl = ""

func handleMessage(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Failed to upgrade connection:", err)
		return
	}

	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}
		log.Printf("Recieved message: %s \n", message)

		data := Coordinate{}
		json.Unmarshal([]byte(message), &data)
		actions[data]()
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please pass your IFTT maker webhook key as the first argument")
		os.Exit(1)
	}
	webhookKey := os.Args[1:][0]
	if webhookKey == "" {
		fmt.Println("Please pass your IFTT maker webhook key as the first argument")
		os.Exit(1)
	}

	webhookUrl = "https://maker.ifttt.com/trigger/%s/with/key/" + webhookKey
	fmt.Println(webhookUrl)

	http.HandleFunc("/", handleMessage)
	log.Fatal(http.ListenAndServe(*flag.String("Stream Deck Websocket Server", "localhost:1337", "Stream Deck Websocket Server"), nil))
}
