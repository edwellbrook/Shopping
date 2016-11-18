package main

import (
	"github.com/yosssi/gmq/mqtt"
	mqtt_client "github.com/yosssi/gmq/mqtt/client"
)

func mqtt_publish(topic string, message []byte) {
	_ = mqttClient.Publish(&mqtt_client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(topic),
		Message:   message,
	})
}
