package main

import (
	"bufio"
	"bytes"
	"log"
	"serial_api"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/tarm/serial"
	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

type SerialResponse struct {
	Data  []byte
	Error error
}

var config *Config
var serialPort *serial.Port
var mqttClient *mqtt_client.Client
var serialChan chan SerialResponse

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

// func setupSerial(config Config) func() error {
// 	return func() (err error) {
// 		serialPort, err = serial.OpenPort(config.Serial)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		return err
// 	}
// }
//
// func connectToSerial(config Config) error {
// 	serialBackoff := backoff.NewExponentialBackOff()
// 	serialBackoff.MaxElapsedTime = time.Minute
//
// 	err := backoff.Retry(setupSerial(config), serialBackoff)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func processSerial(response SerialResponse) {
	if response.Error != nil {
		// there was an error, attempt to reconnect to hopefully resolve it
		go setupSerial()
		return
	}

	resp := serial_api.Parse(response.Data)

	switch resp.Type {
	case serial_api.INIT:
		log.Println("Application starting")

	case serial_api.INFO:
		msg := strings.Join(resp.Args, " ")
		log.Println(msg)

	case serial_api.AUTH:
		if len(resp.Args) == 1 {
			cardId := []byte(resp.Args[0])
			authoriseCard(cardId)
		}

	case serial_api.SCAN:
		if len(resp.Args) == 1 {
			foundBeacon(resp.Args[0])
		}

	default:
		log.Printf("Application sent unhandled message: %s\n", resp.Type)
	}
}

func setupSerial() {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	err := backoff.Retry(func() (err error) {
		println("attempting to open ting")

		serialPort, err = serial.OpenPort(config.Serial)
		if err != nil {
			log.Printf("-> %s", err)
		}
		return err
	}, b)

	if err != nil {
		log.Fatal(err)
	}

	defer serialPort.Close()

	reader := bufio.NewReader(serialPort)

	for {
		line, _, err := reader.ReadLine()
		serialChan <- SerialResponse{line, err}

		if err != nil {
			break
		}
	}
}

func main() {
	var err error

	config, err = appConfig()
	if err != nil {
		log.Fatal(err)
	}

	// err = connectToSerial(*config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer serialPort.Close()

	// serialPort, err = serial.OpenPort(config.Serial)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	serialChan = make(chan SerialResponse)

	go setupSerial()

	for {
		select {
		case response := <-serialChan:
			processSerial(response)
		}
	}
}
