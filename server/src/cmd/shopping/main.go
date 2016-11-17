package main

import (
	"bufio"
	"errors"
	"flag"
	"log"

	"github.com/tarm/serial"
)

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

func main() {
	config, err := appConfig()
	if err != nil {
		log.Fatal(err)
	}

	serialConfig := &serial.Config{
		Name: config.SerialName,
		Baud: config.SerialBaud,
	}

	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	reader := bufio.NewReader(port)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%q", line)
	}
}
