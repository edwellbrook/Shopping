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

	var card string

	err := postgres.QueryRow("SELECT id FROM cards WHERE id = $1", cardId).Scan(&card)

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

func loadShoppingList(cardId string) [12]string {
	var list string
	var res [12]string

	postgres.QueryRow("SELECT list FROM cards WHERE id = $1", cardId).Scan(&list)
	l := strings.Trim(list, "{}")

	for i, v := range strings.Split(l, ",") {
		res[i] = v

		if i == 11 {
			break
		}
	}

	return res
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

	case api.HELP:
		if len(r.Args) == 1 {
			foundBeacon(r.Args[0])
		}

	case api.LIST:
		if len(r.Args) == 1 {
			cardId := hex.EncodeToString([]byte(r.Args[0]))
			list := loadShoppingList(cardId)
			device.SendList(list)
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
