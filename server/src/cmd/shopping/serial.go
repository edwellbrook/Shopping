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

func processSerial(r serial_api.Response) {
	// if there was an error, attempt to reconnect to hopefully resolve it
	if r.Error != nil {
		go setupSerial()
		return
	}

	switch r.Type {
	case serial_api.INIT:
		log.Println("Application starting")

	case serial_api.INFO:
		msg := strings.Join(r.Args, " ")
		log.Println(msg)

	case serial_api.AUTH:
		if len(r.Args) == 1 {
			cardId := []byte(r.Args[0])
			authoriseCard(cardId)
		}

	case serial_api.SCAN:
		if len(r.Args) == 1 {
			foundBeacon(r.Args[0])
		}

	default:
		log.Printf("Application sent unhandled message: %s\n", r.Type)
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
		serialChan <- *serial_api.NewResponse(line, err)

		if err != nil {
			break
		}
	}

	serialPort.Close()
}
