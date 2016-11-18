package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"serial_api"
	"strings"

	"github.com/tarm/serial"
)

var serialPort *serial.Port

type Config struct {
	SerialName string
	SerialBaud int
}

func appConfig() (*Config, error) {
	conf := &Config{}

	name := flag.String("com", "", "COM port for transferring data")
	baud := flag.Int("baud", 9600, "Baud rate for COM port")

	flag.Parse()

	if *name == "" {
		return conf, errors.New("A COM port must be specified")
	}

	conf.SerialName = *name
	conf.SerialBaud = *baud

	return conf, nil
}

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
		if _, err := serialPort.Write([]byte{'1'}); err != nil {
			log.Println("Failed to write auth response")
		} else {
			log.Println("Wrote auth response %d", 1)
		}
	default:
		log.Printf("Application sent unhandled message: %s\n", response.Type)
	}
}

func main() {
	config, err := appConfig()
	if err != nil {
		log.Fatal(err)
	}

	serialConfig := &serial.Config{
		Name: config.SerialName,
		Baud: config.SerialBaud,
	}

	serialPort, err = serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer serialPort.Close()

	reader := bufio.NewReader(serialPort)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		parseInput(line)
	}
}
