package main

import (
	"errors"
	"flag"

	"github.com/tarm/serial"
	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

type Config struct {
	Serial *serial.Config
	MQTT   *mqtt_client.ConnectOptions
}

func appConfig() (*Config, error) {
	conf := &Config{}

	name := flag.String("com", "", "COM port for transferring data")
	baud := flag.Int("baud", 9600, "Baud rate for COM port")
	mqtt := flag.String("mqtt", "127.0.0.1:1883", "Address for MQTT server")

	flag.Parse()

	if *name == "" {
		return conf, errors.New("A COM port must be specified")
	}

	conf.Serial = &serial.Config{
		Name: *name,
		Baud: *baud,
	}

	conf.MQTT = &mqtt_client.ConnectOptions{
		Network:  "tcp",
		Address:  *mqtt,
		ClientID: []byte("shopping-client"),
	}

	return conf, nil
}
