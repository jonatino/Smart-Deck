package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var actions = map[Coordinate]func(){
	Coordinate{0, 1}: lightsToggle,
	Coordinate{0, 2}: speakersToggle,
	Coordinate{0, 3}: fanToggle,
	Coordinate{0, 4}: nil,
	Coordinate{1, 0}: nil,
	Coordinate{1, 1}: nil,
	Coordinate{1, 2}: nil,
	Coordinate{1, 3}: nil,
	Coordinate{1, 4}: nil,
	Coordinate{2, 0}: inputSwitch,
	Coordinate{2, 1}: nil,
	Coordinate{2, 2}: nil,
	Coordinate{2, 3}: nil,
	Coordinate{2, 4}: nil,
}

func lightsToggle()   { webhookRequest("lights_toggle") }
func speakersToggle() { webhookRequest("speakers_toggle") }
func fanToggle()      { webhookRequest("fan_toggle") }

func inputSwitch() { webhookRequest("input_switch") }

func deviceSleeping() { webhookRequest("speaker_outlet_off") }
func deviceWake()     { webhookRequest("speaker_outlet_on") }

func webhookRequest(key string) {
	fmt.Println(fmt.Sprintf(webhookUrl, key))
	response, err := http.Post(fmt.Sprintf(webhookUrl, key), "", nil)
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
