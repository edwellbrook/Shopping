package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"serial"
)

type HelpMessage struct {
	CardId   string `json:"id"`
	BeaconId string `json:"beacon"`
	Location string `json:"location"`
}

func authoriseCard(cardId string) bool {
	log.Printf("Authorising card: %v\n", cardId)

	var card string

	err := postgres.QueryRow("SELECT id FROM cards WHERE id = $1", cardId).Scan(&card)

	if err == sql.ErrNoRows {
		return false
	} else if err == nil {
		return true
	} else {
		log.Printf("Database error: %s\n", err)
		return false
	}
}

func publishHelp(cardId string, beaconId string) {
	msg := &HelpMessage{
		CardId:   cardId,
		BeaconId: beaconId,
		Location: beaconId, // will use psql to translate to location
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	publishMessage(MQTTMessage{"/help", data})
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

func processSerialResponse(device *serial.Device, r *serial.Response) {
	switch r.Type {
	case serial.INFO:
		msg := strings.Join(r.Args, " ")
		log.Println(msg)

	case serial.AUTH:
		if len(r.Args) == 1 {
			cardId := hex.EncodeToString([]byte(r.Args[0]))
			success := authoriseCard(cardId)
			device.Authorise(success)
		}

	case serial.HELP:
		if len(r.Args) == 2 {
			cardId := hex.EncodeToString([]byte(r.Args[0]))
			beaconId := hex.EncodeToString([]byte(r.Args[1]))
			publishHelp(cardId, beaconId)
		}

	case serial.LIST:
		if len(r.Args) == 1 {
			cardId := hex.EncodeToString([]byte(r.Args[0]))
			list := loadShoppingList(cardId)
			device.SendList(list)
		}

	default:
		log.Printf("Application sent unhandled message: %+v\n", r)
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
}
