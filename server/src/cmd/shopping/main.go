package main

import (
	"bytes"
	"log"
	"serial_api"

	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

var config *Config
var mqttClient *mqtt_client.Client
var serialChan chan serial_api.Response

func authoriseCard(cardId []byte) {
	log.Printf("Authorising card: %v\n", cardId)

	var response []byte

	// we should check this against an auth server but for now
	// just check against literal
	if bytes.Equal([]byte{131, 20, 142, 171, 45, 195, 1}, cardId) {
		response = []byte{1}
		mqtt_publish("/auth/accepted", cardId)
	} else {
		response = []byte{0}
		mqtt_publish("/auth/denied", cardId)
	}

	if _, err := serialPort.Write(response); err != nil {
		log.Printf("Failed to write auth response: %s\n", err)
	}
}

func foundBeacon(beaconId string) {
	log.Printf("Found beacon: %s", beaconId)
}

func main() {
	config = mustLoadConfig()

	serialChan = make(chan serial_api.Response)

	go setupSerial()

	for {
		select {
		case response := <-serialChan:
			processSerial(response)
		}
	}
}
