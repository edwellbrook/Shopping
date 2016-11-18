package main

import (
	"bufio"
	"log"
	"serial_api"
	"strings"

	"github.com/tarm/serial"
	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

var serialPort *serial.Port
var mqttClient *mqtt_client.Client

func parseInput(input []byte) {
	response := serial_api.Parse(input)

	switch response.Type {
	case serial_api.INIT:
		log.Println("Application starting")
	case serial_api.EXIT:
		log.Println("Application exiting")
	case serial_api.INFO:
		msg := strings.Join(response.Args, " ")
		log.Println(msg)
	case serial_api.AUTH:
		log.Println("Application asking for auth")
		if _, err := serialPort.Write([]byte("1")); err != nil {
			log.Println("Failed to write auth response")
		} else {
			log.Println("Wrote auth response %d", 1)
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
