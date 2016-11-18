package main

import (
	"bufio"
	"bytes"
	"log"
	"serial_api"
	"strings"

	"github.com/tarm/serial"
	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

var serialPort *serial.Port
var mqttClient *mqtt_client.Client

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

func parseInput(input []byte) {
	response := serial_api.Parse(input)

	switch response.Type {
	case serial_api.INIT:
		log.Println("Application starting")

	case serial_api.INFO:
		msg := strings.Join(response.Args, " ")
		log.Println(msg)

	case serial_api.AUTH:
		if len(response.Args) == 1 {
			cardId := []byte(response.Args[0])
			authoriseCard(cardId)
		}

	default:
		log.Printf("Application sent unhandled message: %s\n", response.Type)
	}
}

func processSerial() {
	reader := bufio.NewReader(serialPort)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		parseInput(line)
	}
}

func main() {
	config, err := appConfig()
	if err != nil {
		log.Fatal(err)
	}

	serialPort, err = serial.OpenPort(config.Serial)
	if err != nil {
		log.Fatal(err)
	}
	defer serialPort.Close()

	mqttClient = mqtt_client.New(&mqtt_client.Options{})
	err = mqttClient.Connect(config.MQTT)
	if err != nil {
		log.Fatal(err)
	}
	defer mqttClient.Terminate()

	processSerial()
}
