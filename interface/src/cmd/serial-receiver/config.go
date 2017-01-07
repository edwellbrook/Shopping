package main

import (
	"errors"
	"flag"
	"log"

	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

var ErrNoCOM = errors.New("A COM port must be specified")

type Config struct {
	Serial   string
	MQTT     *mqtt_client.ConnectOptions
	Postgres string
}

func mustLoadConfig() *Config {
	conf := &Config{}

	name := flag.String("com", "", "COM port for transferring data")
	mqtt := flag.String("mqtt", "127.0.0.1:1883", "Address for MQTT server")
	psql := flag.String("psql", "127.0.0.1:5432", "Address for Postgres server")

	flag.Parse()

	if *name == "" {
		log.Fatal(ErrNoCOM)
	}

	conf.Serial = *name

	conf.MQTT = &mqtt_client.ConnectOptions{
		Network:  "tcp",
		Address:  *mqtt,
		ClientID: []byte("shopping-client"),
	}

	conf.Postgres = *psql

	return conf
}
