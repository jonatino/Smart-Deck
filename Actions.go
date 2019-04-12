package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

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
