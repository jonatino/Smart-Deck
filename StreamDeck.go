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

		if string(message) == "DEVICE-SLEEPING" {
			deviceSleeping()
		} else if string(message) == "DEVICE-WAKE" {
			deviceWake()
		} else {
			data := Coordinate{}
			json.Unmarshal([]byte(message), &data)
			actions[data]()
		}
	}
}

func main() {
	fail := func() {
		fmt.Println("Please pass your Home Assistant base url (e.x http://192.168.0.100:8123 or http://homeassist.example.com) as the first argument")
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fail()
	}

	baseurl := os.Args[1:][0]
	if baseurl == "" {
		fail()
	}

	webhookUrl = baseurl + "/api/webhook/%s"

	fmt.Println(webhookUrl)

	http.HandleFunc("/", handleMessage)
	log.Fatal(http.ListenAndServe(*flag.String("Stream Deck Websocket Server", "localhost:1337", "Stream Deck Websocket Server"), nil))
}
