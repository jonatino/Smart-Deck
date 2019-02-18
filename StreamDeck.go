package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{}

var webhookUrl = ""

var actions = map[Coordinate]func(){
	Coordinate{0, 1}: lightsOn,
	Coordinate{0, 2}: lightsOff,
	Coordinate{0, 3}: fanOn,
	Coordinate{0, 4}: fanOff,
	Coordinate{1, 0}: watchTv,
	Coordinate{1, 1}: netflix,
	Coordinate{1, 2}: amazon,
	Coordinate{1, 3}: viewPc,
	Coordinate{1, 4}: tvOff,
}

func lightsOn()  { webhookRequest("lights_on") }
func lightsOff() { webhookRequest("lights_off") }

func fanOn()  { webhookRequest("fan_on") }
func fanOff() { webhookRequest("fan_off") }

func watchTv() { webhookRequest("tv_on") }
func netflix() { webhookRequest("tv_netflix") }
func amazon()  { webhookRequest("tv_amazon") }
func viewPc()  { webhookRequest("tv_pc") }
func tvOff()   { webhookRequest("tv_off") }

func webhookRequest(key string) {
	response, err := http.Get(fmt.Sprintf(webhookUrl, key))
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
	}
}

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

type Coordinate struct {
	Row    int32 `json:"row"`
	Column int32 `json:"column"`
}
