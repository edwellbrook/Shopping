package serial_device

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"serial_api"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/tarm/serial"
)

type Device struct {
	config *serial.Config
	port   *serial.Port
	reader *bufio.Reader
}

func NewDevice(name string) *Device {
	return &Device{config: &serial.Config{
		Name: name,
		Baud: 9600,
	}}
}

func (s *Device) handshake() (err error) {
	// set up new reader
	s.reader = bufio.NewReader(s.port)

	// generate random bytes
	checkBytes := make([]byte, 5)
	rand.Read(checkBytes)

	// send random bytes
	s.port.Write(checkBytes)

	line, _, err := s.reader.ReadLine()
	if err != nil {
		return err
	}

	if string(line) != fmt.Sprintf("%s%s%s", serial_api.ECHO, serial_api.DELIMITTER, checkBytes) {
		return errors.New("Echo handshake failed")
	}

	s.port.Write([]byte("READY"))

	return nil
}

func (s *Device) Open() (err error) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	return backoff.Retry(func() error {
		s.port, err = serial.OpenPort(s.config)
		if err != nil {
			return err
		}

		return s.handshake()
	}, b)
}

func (s *Device) Read() (*serial_api.Response, error) {
	line, _, err := s.reader.ReadLine()
	if err != nil {
		return nil, err
	}

	return serial_api.NewResponse(line), nil
}

func (s *Device) Authorise(b bool) {
	if b == true {
		s.Write([]byte("AUTH1"))
	} else {
		s.Write([]byte("AUTH0"))
	}
}

func (s *Device) Write(data []byte) {
	s.port.Write(data)
}

func (s *Device) Reset() error {
	s.port.Close()
	return s.Open()
}
