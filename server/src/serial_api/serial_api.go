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
	Type string
	Args []string
}

func Parse(input []byte) *Response {
	str := string(input)
	args := strings.Split(str, DELIMITTER)

	if len(args) > 1 {
		return &Response{Type: args[0], Args: args[1:]}
	}

	return &Response{Type: args[0]}
}
