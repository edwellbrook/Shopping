package main

import (
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

func publishMessage(m MQTTMessage) error {
	return mqttClient.Publish(&mqtt_client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(m.Topic),
		Message:   m.Data,
	})
}

func openMQTT(config *mqtt_client.ConnectOptions) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	return backoff.Retry(func() error {
		mqttClient = mqtt_client.New(nil)

		return mqttClient.Connect(config)
	}, b)
}
