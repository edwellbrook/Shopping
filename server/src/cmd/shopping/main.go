package main

import (
	"bufio"
	"flag"
	"log"

	"github.com/tarm/serial"
)

func main() {
	com := flag.String("com", "", "COM port for transferring data")
	baud := flag.Int("baud", 9600, "Baud rate for COM port")

	flag.Parse()

	if *com == "" {
		log.Fatal("A COM port must be specified")
	}

	sConf := &serial.Config{Name: *com, Baud: *baud}

	port, err := serial.OpenPort(sConf)
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
