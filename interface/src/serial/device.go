package serial

import (
	"bufio"
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

func (s *Device) Open() (err error) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	return backoff.Retry(func() error {
		s.port, err = serial.OpenPort(s.config)
		return err
	}, b)
}

func (s *Device) Read() (*Response, error) {
	if s.reader == nil {
		s.reader = bufio.NewReader(s.port)
	}

	line, _, err := s.reader.ReadLine()
	if err != nil {
		return nil, err
	}

	return newResponse(line), nil
}

func (s *Device) Authorise(b bool) {
	if b == true {
		s.write([]byte("AUTH1"))
	} else {
		s.write([]byte("AUTH0"))
	}
}

func (s *Device) SendList(list [12]string) {
	s.port.Write([]byte("LLOAD"))

	for _, item := range list {
		s.write(append([]byte(item), 0))
	}
}

func (s *Device) write(data []byte) {
	s.port.Write(data)
}
