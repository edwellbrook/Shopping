package main

import (
	"database/sql"
	"encoding/hex"
	"log"
	"strings"
	"sync"

	api "serial_api"
	serial "serial_device"
)

func authoriseCard(cardId string) bool {
	log.Printf("Authorising card: %v\n", cardId)

	var userId int

	err := postgres.QueryRow("SELECT user_id FROM cards WHERE card_id = $1", cardId).Scan(&userId)

	if err == sql.ErrNoRows {
		publishMessage(MQTTMessage{"/auth/denied", []byte(cardId)})
		return false
	} else if err == nil {
		publishMessage(MQTTMessage{"/auth/accepted", []byte(cardId)})
		return true
	} else {
		log.Printf("Database error: %s\n", err)
		return false
	}
}

func foundBeacon(beaconId string) {
	log.Printf("Found beacon: %s", beaconId)
}

func processSerialResponse(device *serial.Device, r *api.Response) {
	switch r.Type {
	case api.INIT:
		log.Println("Application starting")

	case api.INFO:
		msg := strings.Join(r.Args, " ")
		log.Println(msg)

	case api.AUTH:
		if len(r.Args) == 1 {
			cardId := hex.EncodeToString([]byte(r.Args[0]))
			go func() {
				success := authoriseCard(cardId)
				device.Authorise(success)
			}()
		}

	case api.SCAN:
		if len(r.Args) == 1 {
			foundBeacon(r.Args[0])
		}

	default:
		log.Printf("Application sent unhandled message: %s\n", r.Type)
	}
}

func main() {
	conf := mustLoadConfig()
	wait := sync.WaitGroup{}

	wait.Add(1)
	go func() {
		if err := openMQTT(conf.MQTT); err != nil {
			log.Fatal(err)
		}
		wait.Done()
	}()

	wait.Add(1)
	go func() {
		if err := openPostgres(conf.Postgres); err != nil {
			log.Fatal(err)
		}

		if postgres == nil {
			log.Fatal("postgres variable was not set")
		}

		wait.Done()
	}()

	wait.Wait()
	log.Println("Services connected")

	device := serial.NewDevice(conf.Serial)
	device.Open()

	log.Println("Reading serial data")

	for {
		response, err := device.Read()
		if err != nil {
			log.Fatal(err)
		}

		processSerialResponse(device, response)
	}

	// for {
	// 	select {
	// 	case response := <-serialChan:
	// 		log.Println(response)
	// 		processSerial(response)
	// 	case message := <-mqttChan:
	// 		publishMessage(message)
	// 	}
	// }
}
