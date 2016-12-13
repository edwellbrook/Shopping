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

	// delimitter
	DELIMITTER = ":"
)

type Response struct {
	Raw   []byte
	Type  string
	Args  []string
	Error error
}

func NewResponse(raw []byte, err error) *Response {
	r := &Response{
		Raw:   raw,
		Error: err,
	}

	if err == nil {
		r.parseRaw()
	}

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
