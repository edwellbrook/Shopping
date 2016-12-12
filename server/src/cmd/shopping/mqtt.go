package main

import (
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/yosssi/gmq/mqtt"
	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

var mqttClient *mqtt_client.Client

type MQTTMessage struct {
	Topic string
	Data  []byte
}

func publishMessage(m MQTTMessage) {
	err := mqttClient.Publish(&mqtt_client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(m.Topic),
		Message:   m.Data,
	})

	if err != nil {
		log.Println(err)
	}
}

func openMQTT() error {
	mqttClient = mqtt_client.New(&mqtt_client.Options{})

	err := mqttClient.Connect(config.MQTT)
	if err != nil {
		log.Println(err)
	}
	return err
}

func setupMQTT() {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	err := backoff.Retry(openMQTT, b)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT: %s\n", err)
	}
}
