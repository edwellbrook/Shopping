package main

import (
	"bufio"
	"log"
	"serial_api"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/tarm/serial"
)

var serialPort *serial.Port

type SerialResponse struct {
	Data  []byte
	Error error
}

func processSerial(response SerialResponse) {
	// if there was an error, attempt to reconnect to hopefully resolve it
	if response.Error != nil {
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

func openSerial() (err error) {
	serialPort, err = serial.OpenPort(config.Serial)
	if err != nil {
		log.Println(err)
	}
	return err
}

func setupSerial() {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	err := backoff.Retry(openSerial, b)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(serialPort)

	for {
		line, _, err := reader.ReadLine()
		serialChan <- SerialResponse{line, err}

		if err != nil {
			break
		}
	}

	serialPort.Close()
}
