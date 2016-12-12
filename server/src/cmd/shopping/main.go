package main

import (
	"bytes"
	"log"
	"serial_api"
)

var config *Config
var serialChan chan serial_api.Response
var mqttChan chan MQTTMessage

func authoriseCard(cardId []byte) {
	log.Printf("Authorising card: %v\n", cardId)

	var response []byte
	var message MQTTMessage

	// we should check this against an auth server but for now
	// just check against literal
	if bytes.Equal([]byte{131, 20, 142, 171, 45, 195, 1}, cardId) {
		response = []byte{1}
		message = MQTTMessage{"/auth/accepted", cardId}
	} else {
		response = []byte{0}
		message = MQTTMessage{"/auth/denied", cardId}
	}

	mqttChan <- message

	if _, err := serialPort.Write(response); err != nil {
		log.Printf("Failed to write auth response: %s\n", err)
	}
}

func foundBeacon(beaconId string) {
	log.Printf("Found beacon: %s", beaconId)
}

func main() {
	config = mustLoadConfig()

	serialChan = make(chan serial_api.Response, 1)
	mqttChan = make(chan MQTTMessage, 1)

	go setupSerial()
	go setupMQTT()

	for {
		select {
		case response := <-serialChan:
			processSerial(response)
		case message := <-mqttChan:
			publishMessage(message)
		}
	}
}
