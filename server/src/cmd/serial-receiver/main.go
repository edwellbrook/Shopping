package main

import (
	"database/sql"
	"log"
	"serial_api"
)

var config *Config
var serialChan chan serial_api.Response
var mqttChan chan MQTTMessage

func authoriseCard(cardId string) {
	log.Printf("Authorising card: %v\n", cardId)

	var response []byte
	var userId int

	err := postgres.QueryRow("SELECT user_id FROM cards WHERE card_id = $1", cardId).Scan(&userId)

	switch {
	case err == sql.ErrNoRows:
		response = []byte{0}
		mqttChan <- MQTTMessage{"/auth/denied", []byte(cardId)}
	case err != nil:
		response = []byte{0}
		log.Printf("Database error: %s\n", err)
	default:
		response = []byte{1}
		mqttChan <- MQTTMessage{"/auth/accepted", []byte(cardId)}
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

	serialChan = make(chan serial_api.Response, 1)
	mqttChan = make(chan MQTTMessage, 1)

	go setupSerial()
	go setupMQTT()
	go setupDatabase()

	for {
		select {
		case response := <-serialChan:
			processSerial(response)
		case message := <-mqttChan:
			publishMessage(message)
		}
	}
}
