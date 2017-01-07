package serial_api

import (
	"strings"
)

const (
	// system has initialised
	INIT = "INIT"

	// system has requested help
	HELP = "HELP"

	// system has scanned beacon
	SCAN = "SCAN"

	// system is authorising card
	AUTH = "AUTH"

	// system has logged information
	INFO = "INFO"

	// system is echoing handshake bytes
	ECHO = "ECHO"

	// delimitter
	DELIMITTER = ":"
)

type Response struct {
	Raw  []byte
	Type string
	Args []string
}

func NewResponse(raw []byte) *Response {
	r := &Response{Raw: raw}
	r.parseRaw()
	return r
}

func (r *Response) parseRaw() {
	str := string(r.Raw)
	args := strings.Split(str, DELIMITTER)

	r.Type = args[0]

	if len(args) > 1 {
		r.Args = args[1:]
	}
}
